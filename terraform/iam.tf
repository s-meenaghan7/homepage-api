resource "aws_iam_role" "lambda_exec" {
  name        = "lambda-default-execution-role"
  description = "Default execution role for Lambda functions invoked by API Gateway with DynamoDB access."

  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : "sts:AssumeRole",
        "Principal" : {
          "Service" : "lambda.amazonaws.com"
        }
      },
      {
        "Effect" : "Allow",
        "Action" : "sts:AssumeRole",
        "Principal" : {
          "Service" : "apigateway.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_exec" {
  name = "lambda-default-execution-policy"
  role = aws_iam_role.lambda_exec.id

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "CloudWatchLogsAccess",
        "Effect" : "Allow",
        "Action" : [
          "logs:*"
        ],
        "Resource" : "*"
      }
    ]
  })
}

resource "aws_iam_role_policy" "dynamodb_access" {
  name = "homepage-visitor-dynamodb-access"
  role = aws_iam_role.lambda_exec.id

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "VisitorTableAccess",
        "Effect" : "Allow",
        "Action" : [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ],
        "Resource" : aws_dynamodb_table.visitor_counter.arn
      }
    ]
  })
}

# AWS-Managed Policy attachment
resource "aws_iam_role_policy_attachment" "lambda_basic_exec" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "deployment" {
  name        = "HomepageApiDeploymentPolicy"
  description = "Permissions policy to deploy backend resources for homepage."

  policy = file("${path.module}/iam/policies/deployment-policy.json")
}

resource "aws_iam_role" "backend_deploy" {
  name        = "homepage-backend-deploy-gha-oidc"
  description = "This role is used to deploy backend infrastructure and code to AWS from GitHub."

  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Principal" : {
          "Federated" : "arn:aws:iam::659077917555:oidc-provider/token.actions.githubusercontent.com"
        },
        "Action" : "sts:AssumeRoleWithWebIdentity",
        "Condition" : {
          "StringEquals" : {
            "token.actions.githubusercontent.com:aud" : "sts.amazonaws.com"
          },
          "StringLike" : {
            "token.actions.githubusercontent.com:sub" : [
              "repo:s-meenaghan7/homepage-api:*",
              "repo:s-meenaghan7/homepage-api:*"
            ]
          }
        }
      }
    ]
  })
}

# REMOVED 2/18/2026 - This role is currently not used due to complexity setting a secure, but broadly-scoped policy for Terraform deployments.
# resource "aws_iam_role_policy_attachment" "backend_deploy" {
#   role       = aws_iam_role.backend_deploy.name
#   policy_arn = aws_iam_policy.deployment.arn
# }

resource "aws_iam_role_policy_attachment" "terraform_deploy" {
  role       = aws_iam_role.backend_deploy.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}

resource "aws_iam_policy" "deploy_site" {
  name        = "HomepageDeploymentPolicy"
  description = "Allows GitHub Actions to deploy website files to S3 and invalidate CloudFront cache for homepage."

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "S3DeploymentPermissions",
        "Effect" : "Allow",
        "Action" : [
          "s3:PutObject",
          "s3:PutObjectAcl",
          "s3:GetObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ],
        "Resource" : [
          "arn:aws:s3:::seanmeenaghan.com",
          "arn:aws:s3:::seanmeenaghan.com/*"
        ]
      },
      {
        "Sid" : "CloudFrontInvalidationPermissions",
        "Effect" : "Allow",
        "Action" : [
          "cloudfront:CreateInvalidation",
          "cloudfront:GetInvalidation"
        ],
        "Resource" : "arn:aws:cloudfront::659077917555:distribution/E5WTJH2OIG4HG"
      }
    ]
  })
}

resource "aws_iam_role" "frontend_deploy" {
  name        = "gh-oidc"
  description = "Role for GitHub OIDC IdP. FOR FRONTEND RESOURCE CHANGES ONLY."

  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Principal" : {
          "Federated" : "arn:aws:iam::659077917555:oidc-provider/token.actions.githubusercontent.com"
        },
        "Action" : "sts:AssumeRoleWithWebIdentity",
        "Condition" : {
          "StringEquals" : {
            "token.actions.githubusercontent.com:aud" : "sts.amazonaws.com"
          },
          "StringLike" : {
            "token.actions.githubusercontent.com:sub" : [
              "repo:s-meenaghan7/homepage:*",
              "repo:s-meenaghan7/homepage:*"
            ]
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "frontend_deploy" {
  role       = aws_iam_role.frontend_deploy.name
  policy_arn = aws_iam_policy.deploy_site.arn
}

# ------------------------------------------------------------------------------
# IAM — API Gateway must be allowed to write to CloudWatch Logs
# ------------------------------------------------------------------------------
resource "aws_iam_role" "api_gw_cloudwatch" {
  name = "api-gateway-cloudwatch-logs-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "apigateway.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "api_gw_cloudwatch" {
  role       = aws_iam_role.api_gw_cloudwatch.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}

# Registers the IAM role with API Gateway at the account level.
resource "aws_api_gateway_account" "main" {
  cloudwatch_role_arn = aws_iam_role.api_gw_cloudwatch.arn
}

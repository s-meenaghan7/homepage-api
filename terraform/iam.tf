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
        "Resource" : "arn:aws:dynamodb:us-west-2:659077917555:table/homepage-visitor-table-*"
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

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "DynamoDBPermissions",
        "Effect" : "Allow",
        "Action" : [
          "dynamodb:CreateTable",
          "dynamodb:DeleteTable",
          "dynamodb:DescribeTable",
          "dynamodb:DescribeTimeToLive",
          "dynamodb:UpdateTable",
          "dynamodb:UpdateTimeToLive",
          "dynamodb:ListTables",
          "dynamodb:ListTagsOfResource",
          "dynamodb:TagResource",
          "dynamodb:UntagResource",
          "dynamodb:DescribeContinuousBackups",
          "dynamodb:UpdateContinuousBackups",
          "dynamodb:DescribeContributorInsights",
          "dynamodb:UpdateContributorInsights"
        ],
        "Resource" : [
          "arn:aws:dynamodb:us-west-2:659077917555:table/*"
        ]
      },
      {
        "Sid" : "LambdaPermissions",
        "Effect" : "Allow",
        "Action" : [
          "lambda:CreateFunction",
          "lambda:DeleteFunction",
          "lambda:GetFunction",
          "lambda:GetFunctionConfiguration",
          "lambda:UpdateFunctionCode",
          "lambda:UpdateFunctionConfiguration",
          "lambda:ListFunctions",
          "lambda:ListVersionsByFunction",
          "lambda:PublishVersion",
          "lambda:CreateAlias",
          "lambda:DeleteAlias",
          "lambda:GetAlias",
          "lambda:UpdateAlias",
          "lambda:AddPermission",
          "lambda:RemovePermission",
          "lambda:GetPolicy",
          "lambda:TagResource",
          "lambda:UntagResource",
          "lambda:ListTags",
          "lambda:InvokeFunction"
        ],
        "Resource" : [
          "arn:aws:lambda:us-west-2:659077917555:function:*"
        ]
      },
      {
        "Sid" : "APIGatewayPermissions",
        "Effect" : "Allow",
        "Action" : [
          "apigateway:GET",
          "apigateway:POST",
          "apigateway:PUT",
          "apigateway:PATCH",
          "apigateway:DELETE",
          "apigateway:UpdateRestApiPolicy"
        ],
        "Resource" : [
          "arn:aws:apigateway:us-west-2::/restapis/*",
          "arn:aws:apigateway:us-west-2::/restapis"
        ]
      },
      {
        "Sid" : "IAMRolePermissions",
        "Effect" : "Allow",
        "Action" : [
          "iam:CreateRole",
          "iam:DeleteRole",
          "iam:GetRole",
          "iam:AttachRolePolicy",
          "iam:DetachRolePolicy",
          "iam:PutRolePolicy",
          "iam:DeleteRolePolicy",
          "iam:GetRolePolicy",
          "iam:ListRolePolicies",
          "iam:ListAttachedRolePolicies",
          "iam:TagRole",
          "iam:UntagRole",
          "iam:ListInstanceProfilesForRole"
        ],
        "Resource" : [
          "arn:aws:iam::659077917555:role/*"
        ]
      },
      {
        "Sid" : "IAMPassRolePermissions",
        "Effect" : "Allow",
        "Action" : [
          "iam:PassRole"
        ],
        "Resource" : [
          "arn:aws:iam::659077917555:role/*"
        ],
        "Condition" : {
          "StringEquals" : {
            "iam:PassedToService" : [
              "lambda.amazonaws.com",
              "apigateway.amazonaws.com"
            ]
          }
        }
      },
      {
        "Sid" : "IAMPolicyPermissions",
        "Effect" : "Allow",
        "Action" : [
          "iam:CreatePolicy",
          "iam:DeletePolicy",
          "iam:GetPolicy",
          "iam:GetPolicyVersion",
          "iam:ListPolicyVersions",
          "iam:CreatePolicyVersion",
          "iam:DeletePolicyVersion"
        ],
        "Resource" : [
          "arn:aws:iam::659077917555:policy/*"
        ]
      },
      {
        "Sid" : "CloudWatchLogsPermissions",
        "Effect" : "Allow",
        "Action" : [
          "logs:CreateLogGroup",
          "logs:DeleteLogGroup",
          "logs:DescribeLogGroups",
          "logs:ListTagsLogGroup",
          "logs:TagLogGroup",
          "logs:UntagLogGroup",
          "logs:PutRetentionPolicy",
          "logs:DeleteRetentionPolicy"
        ],
        "Resource" : "*"
      }
    ]
  })
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

resource "aws_iam_role_policy_attachment" "backend_deploy" {
  role       = aws_iam_role.backend_deploy.name
  policy_arn = aws_iam_policy.deployment.arn
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

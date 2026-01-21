resource "aws_apigatewayv2_api" "homepage_api" {
  name                         = "homepage"
  protocol_type                = "HTTP"
  description                  = "homepage API gateway"
  disable_execute_api_endpoint = true
  route_selection_expression   = "$request.method $request.path"

  cors_configuration {
    max_age       = 3600
    allow_methods = ["GET", "POST"]
    allow_origins = ["https://www.seanmeenaghan.com", "https://seanmeenaghan.com"]
  }
}

resource "aws_apigatewayv2_stage" "production" {
  api_id      = aws_apigatewayv2_api.homepage_api.id
  name        = "production"
  auto_deploy = true
}

resource "aws_apigatewayv2_route" "get_visitor" {
  api_id    = aws_apigatewayv2_api.homepage_api.id
  route_key = "GET /visitor"
  target    = "integrations/${aws_apigatewayv2_integration.get_visitor.id}"
}

resource "aws_apigatewayv2_route" "post_visitor" {
  api_id    = aws_apigatewayv2_api.homepage_api.id
  route_key = "POST /visitor"
  target    = "integrations/${aws_apigatewayv2_integration.post_visitor.id}"
}

resource "aws_lambda_function" "visitor_api" {
  function_name = "visitor-api"
  role          = "arn:aws:iam::659077917555:role/lambda-default-execution-role"
  runtime       = "provided.al2023"
  handler       = "bootstrap"

  # this file must exist, but won't be used for updates
  filename = "placeholder.zip"

  # Ignore code-related changes; GitHub Actions handles deployments
  lifecycle {
    ignore_changes = [
      filename
    ]
  }
}

resource "aws_apigatewayv2_integration" "get_visitor" {
  api_id           = aws_apigatewayv2_api.homepage_api.id
  integration_type = "AWS_PROXY"

  connection_type    = "INTERNET"
  description        = "GET /visitor Lambda integration"
  integration_method = "GET"
  integration_uri    = aws_lambda_function.visitor_api.arn
}

resource "aws_apigatewayv2_integration" "post_visitor" {
  api_id           = aws_apigatewayv2_api.homepage_api.id
  integration_type = "AWS_PROXY"

  connection_type    = "INTERNET"
  description        = "POST /visitor Lambda integration"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.visitor_api.arn
}

resource "aws_apigatewayv2_domain_name" "api_domain_name" {
  domain_name = "api.seanmeenaghan.com"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.api_cert.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_acm_certificate" "api_cert" {
  domain_name       = "api.seanmeenaghan.com"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_apigatewayv2_api_mapping" "api_mapping" {
  api_id          = aws_apigatewayv2_api.homepage_api.id
  domain_name     = aws_apigatewayv2_domain_name.api_domain_name.id
  stage           = "production"
  api_mapping_key = "v1/visitor"
}

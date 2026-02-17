# ==============================================================================
# CloudWatch Log Group for HTTP API Gateway Access Logging
# ==============================================================================

resource "aws_cloudwatch_log_group" "api_access_logs" {
  name              = "/aws/apigateway/${aws_apigatewayv2_api.homepage_api.name}/access-logs"
  retention_in_days = 30

  tags = {
    ManagedBy = "Terraform"
  }
}

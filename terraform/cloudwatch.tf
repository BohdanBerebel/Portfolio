resource "aws_cloudwatch_log_group" "app" {
  name              = "/portfolio/app"
  retention_in_days = 14

  tags = {
    Name = "${local.project_name}-app-logs"
  }
}
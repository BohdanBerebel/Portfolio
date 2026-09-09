resource "aws_iam_role" "github_deploy" {
  name = "${local.project_name}-github-deploy"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"

    Statement = [
      {
        Effect = "Allow"

        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }

        Action = "sts:AssumeRoleWithWebIdentity"

        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }

          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:157813971/Portfolio@1358335905:ref:refs/heads/main"
          }
        }
      }
    ]
  })

  tags = {
    Name = "${local.project_name}-github-deploy"
  }
}

resource "aws_iam_role_policy" "github_deploy_ssm" {
  name = "${local.project_name}-github-deploy-ssm"
  role = aws_iam_role.github_deploy.id

  policy = jsonencode({
    Version = "2012-10-17"

    Statement = [
      {
        Effect = "Allow"

        Action = [
          "ssm:SendCommand"
        ]

        Resource = [
          "arn:aws:ec2:eu-central-1:*:instance/${aws_instance.app.id}",
          "arn:aws:ssm:eu-central-1::document/AWS-RunShellScript"
        ]
      },
      {
        Effect = "Allow"

        Action = [
          "ssm:GetCommandInvocation"
        ]

        Resource = "*"
      }
    ]
  })
}
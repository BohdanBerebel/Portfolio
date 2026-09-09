data "aws_ssm_parameter" "amazon_linux" {
  name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64"
}

resource "aws_key_pair" "app" {
  key_name   = "${local.project_name}-app"
  public_key = file(pathexpand("~/.ssh/portfolio-aws.pub"))
}

resource "aws_instance" "app" {
  ami           = data.aws_ssm_parameter.amazon_linux.value
  instance_type = "t3.micro"
  key_name      = aws_key_pair.app.key_name

  subnet_id                   = aws_subnet.public.id
  vpc_security_group_ids      = [aws_security_group.ec2.id]
  associate_public_ip_address = true

  root_block_device {
    volume_size = 30
    volume_type = "gp3"
    encrypted   = true
  }

  user_data = <<-EOF
    #!/bin/bash

    set -e

    dnf update -y
    dnf install -y docker git

    systemctl enable --now docker

    usermod -aG docker ec2-user

    mkdir -p /opt/portfolio
    chown ec2-user:ec2-user /opt/portfolio
  EOF

  tags = {
    Name = "${local.project_name}-app"
  }
}
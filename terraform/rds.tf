resource "aws_db_subnet_group" "postgres" {
  name = "${local.project_name}-postgres"

  subnet_ids = [
    aws_subnet.private_1a.id,
    aws_subnet.private_1b.id
  ]

  tags = {
    Name = "${local.project_name}-postgres"
  }
}

resource "aws_db_instance" "postgres" {
  identifier = "${local.project_name}-postgres"

  engine         = "postgres"
  engine_version = "17"

  instance_class      = "db.t3.micro"
  allocated_storage   = 20
  storage_type        = "gp3"
  storage_encrypted   = true
  publicly_accessible = false
  skip_final_snapshot = true
  deletion_protection = false

  db_name  = "notes_app"
  username = "postgres"
  password = var.db_password

  port = 5432

  db_subnet_group_name   = aws_db_subnet_group.postgres.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  backup_retention_period = 1

  tags = {
    Name = "${local.project_name}-postgres"
  }
}
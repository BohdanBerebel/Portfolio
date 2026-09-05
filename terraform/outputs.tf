output "ec2_public_ip" {
  description = "Public IPv4 address of the application EC2 instance"
  value       = aws_instance.app.public_ip
}
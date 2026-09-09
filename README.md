# Portfolio

Full-stack portfolio application built with React, Go, Python, PostgreSQL, Docker, and AWS.

The project includes authentication, notes, real-time chat, containerized services, infrastructure as code, and automated deployment.

## Tech Stack

- **Frontend:** React, TypeScript, Vite, Nginx
- **Backend:** Go, Gin, Python, FastAPI
- **Database:** PostgreSQL
- **Containers:** Docker, Docker Compose
- **Infrastructure:** AWS, Terraform
- **CI/CD:** GitHub Actions, GitHub Container Registry
- **Deployment:** AWS OIDC, Systems Manager

## Architecture

The application consists of:

- React frontend served by Nginx
- Go API for authentication and notes
- Python API for chat and WebSockets
- PostgreSQL database

The backend services communicate with the same PostgreSQL database.

```text
                         Internet
                            │
                            ▼
                    ┌───────────────┐
                    │      EC2      │
                    │    Nginx      │
                    └───────┬───────┘
                            │
                 ┌──────────┴──────────┐
                 ▼                     ▼
          ┌─────────────┐       ┌─────────────┐
          │   Go API    │       │ Python API  │
          │    :8080    │       │    :8000    │
          └──────┬──────┘       └──────┬──────┘
                 │                     │
                 └──────────┬──────────┘
                            ▼
                    ┌───────────────┐
                    │ RDS PostgreSQL│
                    └───────────────┘
```

## Local Development

Requirements:

- Docker
- Docker Compose

Start the application:

```bash
docker compose up --build
```

The application is available at:

```text
http://localhost:3000
```

## Production

The application is deployed to AWS using Terraform and Docker Compose.

Production infrastructure includes:

- Amazon EC2
- Amazon RDS PostgreSQL
- VPC and Security Groups
- IAM
- GitHub OIDC

PostgreSQL runs on a private RDS instance, while the application services run as Docker containers on EC2.

## CI/CD

Deployment is automated with GitHub Actions.

```text
Git push
   │
   ▼
GitHub Actions
   │
   ▼
Build Docker images
   │
   ▼
GitHub Container Registry
   │
   ▼
AWS OIDC
   │
   ▼
AWS Systems Manager
   │
   ▼
EC2
   │
   ▼
Docker Compose
```

GitHub Actions authenticates to AWS using OIDC instead of long-lived AWS access keys.

AWS Systems Manager is used to execute deployment commands on EC2.

## Security

- RDS is private and not publicly accessible
- Production secrets are stored in AWS Parameter Store
- GitHub Actions authenticates with AWS through OIDC
- Deployment uses AWS Systems Manager
- Passwords are hashed with bcrypt
- APIs use JWT authentication

## Testing

Go:

```bash
cd go_app
go test ./...
```

Python:

```bash
cd python_app
uv run pytest
```

## Future Improvements

- HTTPS and custom domain
- CloudWatch monitoring
- Additional automated tests
- Further production hardening

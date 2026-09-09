#!/bin/bash

set -euo pipefail

cd /opt/portfolio

REGION="eu-central-1"
ENV_FILE="/opt/portfolio/.env"

trap 'rm -f "$ENV_FILE"' EXIT

DB_HOST="$(aws ssm get-parameter \
  --name "/portfolio/DB_HOST" \
  --with-decryption \
  --region "$REGION" \
  --query "Parameter.Value" \
  --output text)"

DB_PORT="$(aws ssm get-parameter \
  --name "/portfolio/DB_PORT" \
  --with-decryption \
  --region "$REGION" \
  --query "Parameter.Value" \
  --output text)"

DB_NAME="$(aws ssm get-parameter \
  --name "/portfolio/DB_NAME" \
  --with-decryption \
  --region "$REGION" \
  --query "Parameter.Value" \
  --output text)"

DB_USER="$(aws ssm get-parameter \
  --name "/portfolio/DB_USER" \
  --with-decryption \
  --region "$REGION" \
  --query "Parameter.Value" \
  --output text)"

DB_PASSWORD="$(aws ssm get-parameter \
  --name "/portfolio/DB_PASSWORD" \
  --with-decryption \
  --region "$REGION" \
  --query "Parameter.Value" \
  --output text)"

JWT_SECRET="$(aws ssm get-parameter \
  --name "/portfolio/JWT_SECRET" \
  --with-decryption \
  --region "$REGION" \
  --query "Parameter.Value" \
  --output text)"

umask 077

cat > "$ENV_FILE" <<EOF
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_NAME=$DB_NAME
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
JWT_SECRET=$JWT_SECRET
EOF

docker compose -f docker-compose.prod.yml pull

docker compose -f docker-compose.prod.yml up -d
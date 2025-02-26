#!/bin/bash

DB_HOST=$DB_HOST
DB_USER=$DB_USER
REGION=$REGION
DB_PORT=$DB_PORT
DB_NAME=$DB_NAME

TOKEN=$(aws rds generate-db-auth-token --hostname $DB_HOST --port $DB_PORT --region $REGION --username $DB_USER)

export DATABASE_URL="mysql://${DB_USER}:${TOKEN}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?tls=skip-verify"

migrate -path=/db/migrations -database="${DATABASE_URL}" up

#!/usr/bin/env bash

# run the rest api development server
# usage example: bin/devserver-rest.sh

# load environment variables from .env file
set -a; [ -f .env ] && . .env; set +a

# run the dev server
go run app/rest/main.go

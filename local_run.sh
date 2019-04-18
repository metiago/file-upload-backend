#!/bin/bash

export HOST=0.0.0.0
export PORT=5000

export DB_HOST=172.26.0.3
export DB_PORT=5432
export DB_USERNAME=postgres
export DB_PASSWORD=12345678
export DB_DATABASE=postgres

export PRI_RSA=rsa/app.rsa
export PUB_RSA=rsa/app.rsa.pub

go run main.go

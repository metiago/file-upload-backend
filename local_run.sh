#!/bin/bash

export HOST=0.0.0.0
export PORT=5000

export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_USERNAME=tiago
export DB_PASSWORD=zero
export DB_DATABASE=zbx1

export PRI_RSA=rsa/app.rsa
export PUB_RSA=rsa/app.rsa.pub

go build -o zbx1

./zbx1

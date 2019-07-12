#!/bin/bash
dos2unix deploy.sh
sleep 1

docker stop $(docker ps -aq)
docker rmi -f $(docker images -f "dangling=true" -q)
docker system prune -a -f

export PASSWORD=12345678
export USER=zbx1
export DATABASE=zbx1

export HOST=0.0.0.0
export PORT=5000

export DB_HOST=172.28.0.2
export DB_PORT=5432
export DB_USERNAME=tiago
export DB_PASSWORD=zero
export DB_DATABASE=zbx1

export PRI_RSA=rsa/app.rsa
export PUB_RSA=rsa/app.rsa.pub

docker-compose up -d --build --force-recreate

go build
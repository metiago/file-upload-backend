#!/bin/bash

export PASSWORD=12345678
export USER=zbx1
export DATABASE=zbx1

# APP
export HOST=0.0.0.0
export PORT=5000

# DB
export DB_HOST=postgres
export DB_PORT=5432 
export DB_USERNAME=postgres
export DB_PASSWORD=12345678
export DB_DATABASE=zbx1

# RSA
export PRI_RSA=/usr/local/rsa/app.rsa
export PUB_RSA=/usr/local/rsa/app.rsa.pub

docker-compose down

sudo chmod 777 -R volumes/

sleep 2

docker stop $(docker ps -aq)

sleep 2

docker rmi $(docker images -f "dangling=true" -q)

sleep 2

sudo docker rm $(docker ps -a -q)

sleep 2

docker-compose up -d --build --force-recreate

sleep 2 

sudo chmod 777 -R volumes/

cp -r $PWD/rsa/ $PWD/volumes/zbx1

sleep 2

docker-compose ps

sleep 2

docker-compose logs -f zbx1
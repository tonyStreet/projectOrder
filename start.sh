#!/bin/bash
docker pull mysql/mysql-server:5.7
docker network create --subnet=192.168.0.0/16 project-order-backend
docker run --name order-mysql --net backend --ip 192.168.0.2 -e MYSQL_ROOT_PASSWORD=my-secret-pw -e MYSQL_DATABASE=logistics -e MYSQL_USER=order-app -e MYSQL_PASSWORD=5sEjLqbLxs -p 3306:3306 -d mysql/mysql-server:5.7
env GOOS=linux GOARCH=amd64 go build -v -o project-order
docker build -t project-order:1.0 .
docker run --name project-order --net backend --ip 192.168.0.3 --dns 8.8.8.8 -p 8080:8080 -d project-order:1.0
docker exec -i order-mysql mysql -uorder-app -p5sEjLqbLxs logistics < create-order-table.sql > start.log

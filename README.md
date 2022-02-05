# Cosmoscan API (backend)

Website: https://cosmoscan.net
Frontend repo: https://github.com/everstake/cosmoscan-front

Cosmoscan is the first data and statistics explorer for the Cosmos network. It provides information oт the overall network operations, governance details, validators and much more. This is still an MVP, so if you have any suggestions, please reach out.

Dependency:
 - Clickhouse
 - Mysql
 - Cosmos node
 - Golang

## How to run ?
At first you need to configure the config.json file.
```sh
cp config.example.json config.json
```
Next step you need to build and run application.
#### Docker-compose way:
```sh
cp docker-compose.example.yml docker-compose.yml
cp docker/.env.example .env
cp docker/clickhouse-users.xml.example docker/clickhouse-users.xml
```
> don`t forget set your passwords
```sh
docker-compose build && docker-compose up -d
```
#### Native way:
> at first setup your dependency and set passwords
```sh
go build && ./cosmoscan-api
```
## Installation:

Installation of backend:

1) Under /root:
git clone https://github.com/IDEP-network/chadscan-back.git
git clone https://github.com/persistenceOne/persistenceCore.git

2) Under /root/chadscan-back/:
cp config.example.json config.json

3) Update /root/chadscan-back/config.json:

{
  "api": {
    "port": "5001",
    "allowed_hosts": [
       "localhost:5001",
       "localhost:3000",
       "165.232.173.187:3000",
       "165.232.173.187:5001"
    ]
  },
  "mysql": {
    "host": "localhost",
    "port": "3306",
    "db": "mytestdb",
    "user": "test",
    "password": "newpassword"
  },
  "clickhouse": {
    "protocol": "http",
    "host": "127.0.0.1",
    "port": 8123,
    "user":  "default",
    "password": "snailison50Meth",
    "database": "cosmoshub3"
  },
  "parser": {
    "node": "http://localhost:1317",
    "batch": 500,
    "fetchers": 5
  },
  "cmc_key": ""
}

5) Create new file .env under /root/chadscan-back/:

DB_NAME=mytestdb
DB_USER=test
DB_PASSWORD=newpassword
DB_HOST=localhost

6) Install:
mysql (Ver 8.0.27-0ubuntu0.20.04.1, ref: https://www.digitalocean.com/community/tutorials/how-to-install-mysql-on-ubuntu-20-04 )
clickhouse

run security script ->

Y, 0, rootpass, Y, Y, Y, Y, Y, Y

sudo mysql

SHOW VARIABLES LIKE 'validate_password%';
SET GLOBAL validate_password.policy = LOW;
CREATE USER 'test'@'localhost' IDENTIFIED BY 'newpassword';
GRANT ALL PRIVILEGES ON *.* TO 'test'@'localhost';
exit

mysql -u test -p => newpassword

CREATE DATABASE mytestdb;
USE mytestdb;
exit

7) Installing clickhouse

https://websiteforstudents.com/how-to-install-and-configure-clickhouse-on-ubuntu-16-04-18-04/ => "default" user password = snailison50Meth

after installation:

- clickhouse-client --password => snailison50Meth
- create database cosmoshub3

8) Install go:

-version: go1.17.2 
https://github.com/golang/go/wiki/Ubuntu

9) persistenceCore:

- persistenceCore/application/constants.go:

//Bech32MainPrefix = "persistence"
Bech32MainPrefix = "idep"


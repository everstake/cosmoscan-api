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

# Hoppscotch Backend API

## Requirements

- MySQL

## Get requirements up and running

### MySQL:

```
docker run \
--name hoppscotch_api_mysql \
-p 127.0.0.1:3306:3306 \
-e MYSQL_ROOT_PASSWORD=hoppscotch \
-e MYSQL_DATABASE=hoppscotch \
-e MYSQL_USER=hoppscotch \
-e MYSQL_PASSWORD=hoppscotch \
-d mysql:8.0
```

### Next runs:
```
docker start hoppscotch_api_mysql
```

## Quickstart

- Copy the config.example.yaml to config.yaml
- Start the API by running `go run main.go`

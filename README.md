# Hoppscotch Backend API

This repository contains an open-source implementation of the Hoppscotch Backend to allow the collaborative features to work on a self-hosted instance of Hoppscotch.

This API has the exact same GraphQL schema as the "official" API.

This API does not store its data in Firebase (which the official probably does), but in a local MySQL database.

## Requirements

- MySQL
- An SMTP mail server

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

## Firebase

You will need to create a Firebase project to get this whole thing running (frontend and backend).

Copy the .env.example in the frontend project to .env en fill in your Firebase credentials.

Generate a [Firebase Admin SDK service account](https://console.firebase.google.com/project/_/settings/serviceaccounts/adminsdk) and reference the JSON from the config.yaml.

## Quickstart

- Copy the config.example.yaml to config.yaml
- Start the API by running `go run main.go`

## Deployment

This backend is available as a [docker image](https://hub.docker.com/repository/docker/jerbob92/hoppscotch-backend/general) `jerbob92/hoppscotch-backend`.

The configuration is expected in the working directory or the folder `/etc/api-config`.

When using docker, the easiest way is to mount a local configuration folder as `/etc/api-config` that contains your `config.yaml` and your Firebase Admin SDK Service User json.

If you're behind a reverse proxy, it might be useful to use `/graphql` for the normal GraphQL traffic, and use `/graphql/ws` for the Subscription/WebSocket traffic.

## Frontend deployment

The default frontend requires some minor changes to connect to your backend since it's not made to connect to a custom backend.

You can see the required changes in [this commit](https://github.com/jerbob92/hoppscotch/commit/2c98d6f8471691156c99a4a3bdd37ab95286fb21). You can also use that branch to do your own deployment.

Be aware that the backend needs to be running to build the frontend because it fetches the schema on build time.


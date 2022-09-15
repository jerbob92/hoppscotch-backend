# Hoppscotch Backend API

This repository contains an open-source implementation of the Hoppscotch Backend to allow the collaborative features to
work on a self-hosted instance of Hoppscotch.

This API has the exact same GraphQL schema as the "official" API.

This API does not store its data in Firebase (which the official probably does), but in a local MySQL database.

## Requirements

- MySQL/Postgres
- An SMTP mail server
- A Firebase project & webapp credentials & Admin SDK credentials

## Get requirements up and running

### MySQL (optional when using `docker-compose`):

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

Generate
a [Firebase Admin SDK service account](https://console.firebase.google.com/project/_/settings/serviceaccounts/adminsdk)
and reference the JSON from the config.yaml.

Create Firestore Database

Go to [Firestore Rules](https://github.com/hoppscotch/hoppscotch/blob/main/firestore.rules) and configure them in your
firestore database.

## Quickstart

- Copy the config.example.yaml to config.yaml
- Start the API by running `go run main.go`

## Quickstart (Docker Compose)

- Copy config.example.yaml to tmp/config.yaml
- Put `Firebase Admin SDK service account` file in tmp folder
- Ensure file mappings at volumes are correct in docker-compose.yml
- run `docker compose up -d` or `docker-compose up -d`

## Deployment

This backend is available as
a [docker image](https://hub.docker.com/r/jerbob92/hoppscotch-backend) `jerbob92/hoppscotch-backend`.

The configuration is expected in the working directory or the folder `/etc/api-config`.

When using docker, the easiest way is to mount a local configuration folder as `/etc/api-config` that contains
your `config.yaml` and your Firebase Admin SDK Service User json.

If you're behind a reverse proxy, it might be useful to use `/graphql` for the normal GraphQL traffic, and
use `/graphql/ws` for the Subscription/WebSocket traffic.

## Frontend deployment

To connect to your own backend, you will need to set the `API_URL` and `API_WS_URL` to the correct URLs for your backend in `packages/hoppscotch-app/.env` when building the frontend. 

There is currently one minor bug in the latest version of Hoppscotch, it has the share URL hardcoded, I have fixed that in my own [fork](https://github.com/jerbob92/hoppscotch/tree/feature/local-hosting), but it's not always needed so I leave it up to you which version you use.

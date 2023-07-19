# stemstr relay

The nostr relay for https://stemstr.app

## Running locally

Run a local Postgres container with [Docker](https://docs.docker.com/compose/install/)

```
docker-compose up postgres
```

Run the relay with Go

```
make build run
```

You now have the Stemstr relay running on `ws://localhost:9000`

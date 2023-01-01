# Search.MAATester.com

A search engine for math contest problems. Deployed at
[search.maatester.com](https://search.maatester.com)

# Setup

## Dependencies

You need to have nodejs, yarn, go, and docker installed on your system. On the
other hand, if you're a nix user, you can just use the provided `flake.nix` and
get everything going with a simple `nix develop` (except docker).

## Package Dependencies

Run `cd ./frontend && yarn`

## Development Servers

1. Start Redis Stack by running either `docker-compose up -d` or `docker compose up -d`.
2. Start the backend server. In the project root, run `go run ./backend/main.go server`.
3. Start the frontend vite server by running `cd ./frontend && yarn dev`.

## Loading AoPS Data

Run `go run ./backend/main.go load && go run ./backend/main.go load -F`; this
will populate the redis database with problems scraped from AoPS.

# Production

Use the included `./docker-compose.example.yml` as an example for how to set up
the containers. There needs to be an environment variable connecting the app to
the redis server.

# Search.MAATester.com

A search engine for math contest problems. Deployed at
[search.maatester.com](https://search.maatester.com)

# Setup

## Dependencies

You need to have nodejs, yarn, go, and docker installed on your system. On the
other hand, if you're a nix user, you can just use the provided `flake.nix` and
get everything going with a simple `nix develop` (except docker, which you
should have on your system).

## Package Dependencies

Run `cd ./frontend && yarn`

## Development Servers

1. Start Meilisearch by running either `docker-compose up -d` or `docker compose up -d`.
2. Start the backend server. In the project root, run `go run ./backend/main.go server`.
3. Start the frontend vite server by running `cd ./frontend && yarn dev`.

## Fetch Data

Run `go run ./backend/main.go dump -C contests/(preference.json)
(output json)`, which will pick a specific set of contests to fetch.

On the other hand, if you are not working with the scraper, you can always
download it from the releases page.

## Loading

Run `go run ./backend/main.go load (json file)`; this
will populate the redis database with problems from the json file invoked.

# Production

Use the included `./docker-compose.example.yml`. The production docker container
is monolithic and includes all the dependencies that it requires.

You must specify the runtime environment variable `SOURCES` with a set of
space-separated URL's that dictate where to download the dataset files.

# The Dataset

The dataset that search.maatester.com uses is updated weekly by GitHub
Actions. It can be accessed at
https://github.com/MAA-Contest-Tester/search/releases/download/dataset/main.json

`main.json` contains a list of 17,000 problems from various short-answer and
olympiad contests. Each entry contains the following fields:

- The source of the problem and its number (i.e. 2023 IMO 1),
- A link to its source on C&P (e.g.
  https://artofproblemsolving.com/community/c3381519)
- A link to its solution/discussion on C&P (e.g.
  https://artofproblemsolving.com/community/c6h3106752p28097575)
- Its problem statement
- Any associated tags (e.g. "combinatorics," "functional equations,"
  inequalities," ...)

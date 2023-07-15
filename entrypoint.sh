#!/bin/sh

/meilisearch &
while ! wget --no-verbose --spider http://localhost:7700; do
  sleep 1s
done
/app/psearch server -L /data/forum.json -D /app/dist

#!/bin/sh


/meilisearch &
sleep 5s
/app/psearch server -L /data/forum.json -D /app/dist

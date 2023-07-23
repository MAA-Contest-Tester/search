#!/bin/sh

/meilisearch &
while ! wget --no-verbose --spider http://localhost:7700; do
  sleep 1s
done
FILE=1
if test -z "$SOURCES"; then
  SOURCES=https://github.com/MAA-Contest-Tester/search/releases/download/dataset/main.json
fi
for url in $SOURCES; do
  mkdir -p "/data" || exit 1
  rm -rf "/data/$FILE.json"
  wget "$url" -O "/data/$FILE.json" || exit 1
  FILE=$((FILE+1))
done

/app/psearch server -L /data/*.json -M /data/*.json -D /app/dist

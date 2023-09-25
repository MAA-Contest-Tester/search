#!/bin/sh

set -x

/meilisearch &
while ! wget --no-verbose --spider http://localhost:7700; do
  sleep 1s
done
FILE=1
if test -z "$SOURCES"; then
  SOURCES="https://github.com/MAA-Contest-Tester/search/releases/download/dataset/computational.json https://github.com/MAA-Contest-Tester/search/releases/download/dataset/nationaloly.json https://github.com/MAA-Contest-Tester/search/releases/download/dataset/international.json"
fi
LARGS=""
MARGS=""
for url in $(echo $SOURCES); do
  mkdir -p "/data" || exit 1
  rm -rf "/data/$FILE.json"
  wget "$url" -O "/data/$FILE.json" || exit 1
  LARGS="$LARGS -L /data/$FILE.json"
  MARGS="$MARGS -M /data/$FILE.json"
  FILE=$((FILE+1))
done

/app/psearch server $LARGS $MARGS -D /app/dist

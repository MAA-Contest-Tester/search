#!/bin/sh


redis-server --loadmodule /opt/redis-stack/lib/redisearch.so --loadmodule /opt/redis-stack/lib/rejson.so --port 6379 &
sleep 5s
/app/psearch server -L /data/forum.json -D /app/dist

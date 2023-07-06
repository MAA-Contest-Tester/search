FROM golang:latest as backend
WORKDIR /build
COPY go.mod go.sum /build/
RUN go mod download
COPY . /build/
RUN make
RUN mkdir -p /data && wget https://github.com/MAA-Contest-Tester/search/releases/download/dataset/data.json -O /data/forum.json

FROM node:alpine as frontend
COPY --from=backend /build/frontend /build/frontend/
WORKDIR /build/frontend
RUN yarn install
RUN yarn build

FROM redis/redis-stack-server:7.2.0-RC2 as redis-stack

FROM redis:7.2-rc-bookworm
WORKDIR /app/

COPY --from=backend /build/out/psearch /app/psearch
COPY --from=backend /data/forum.json /data/forum.json
COPY --from=frontend /build/frontend/dist /app/dist
COPY entrypoint.sh /app/entrypoint.sh

COPY --from=redis-stack /opt/redis-stack/lib/redisearch.so /opt/redis-stack/lib/redisearch.so
COPY --from=redis-stack /opt/redis-stack/lib/rejson.so /opt/redis-stack/lib/rejson.so

CMD ["/app/entrypoint.sh"]

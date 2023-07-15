FROM golang:latest as backend
WORKDIR /build
COPY go.mod go.sum /build/
RUN go mod download
COPY . /build/
RUN make
RUN mkdir -p /data && wget https://github.com/MAA-Contest-Tester/search/releases/download/dataset/main.json -O /data/forum.json

FROM node:alpine as frontend
COPY --from=backend /build/frontend /build/frontend/
WORKDIR /build/frontend
RUN yarn install
RUN yarn build

FROM getmeili/meilisearch:latest
WORKDIR /app/

COPY --from=backend /build/out/psearch /app/psearch
COPY --from=backend /data/forum.json /data/forum.json
COPY --from=frontend /build/frontend/dist /app/dist
COPY entrypoint.sh /app/entrypoint.sh

CMD ["/app/entrypoint.sh"]

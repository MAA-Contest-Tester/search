FROM golang:latest as backend
WORKDIR /build
COPY go.mod go.sum /build/
RUN go mod download
COPY . /build/
RUN make
RUN mkdir -p /data

FROM node:alpine as frontend
WORKDIR /build
COPY --from=backend /build/yarn.lock /build/package.json ./
RUN yarn install
COPY --from=backend /build .
RUN yarn build

FROM getmeili/meilisearch:latest
WORKDIR /app/

COPY --from=backend /build/out/psearch /app/psearch
COPY --from=frontend /build/frontend/dist /app/dist
COPY entrypoint.sh /app/entrypoint.sh

CMD ["/app/entrypoint.sh"]

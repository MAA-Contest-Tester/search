FROM golang:alpine as backend
WORKDIR /build
COPY go.mod go.sum /build/
RUN go mod download
COPY . /build/
RUN go build -o /build/out/psearch

FROM node:alpine as frontend
COPY --from=backend /build/frontend /build/frontend/
WORKDIR /build/frontend
RUN yarn install
RUN yarn build

FROM alpine:latest
WORKDIR /app/
RUN apk add musl
COPY --from=backend /build/entrypoint.sh /entrypoint.sh
COPY --from=backend /build/out/psearch /app/psearch
COPY --from=frontend /build/frontend/dist /app/dist
CMD ["/entrypoint.sh"]

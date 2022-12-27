FROM golang:latest as backend
WORKDIR /build
COPY go.mod go.sum /build/
RUN go mod download
COPY . /build/
RUN make

FROM node:alpine as frontend
COPY --from=backend /build/frontend /build/frontend/
WORKDIR /build/frontend
RUN yarn install
RUN yarn build

FROM alpine:latest
WORKDIR /app/
COPY --from=backend /build/entrypoint.sh /entrypoint.sh
COPY --from=backend /build/out/psearch /app/psearch
COPY --from=frontend /build/frontend/dist /app/dist
RUN /app/psearch dump /data/problems.json
CMD ["/entrypoint.sh"]

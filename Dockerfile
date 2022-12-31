FROM golang:latest as backend
WORKDIR /build
COPY go.mod go.sum /build/
RUN go mod download
COPY . /build/
RUN make
RUN /build/out/psearch dump /data/wiki.json
RUN /build/out/psearch dump -F /data/forum.json

FROM node:alpine as frontend
COPY --from=backend /build/frontend /build/frontend/
WORKDIR /build/frontend
RUN yarn install
RUN yarn build

FROM scratch
WORKDIR /app/
COPY --from=backend /build/out/psearch /app/psearch
COPY --from=backend /data/problems.json /data/wiki.json
COPY --from=frontend /build/frontend/dist /app/dist
CMD ["/app/psearch", "server", "-L", "/data/wiki.json", "-L", "/data/forum.json", "-D", "/app/dist"]

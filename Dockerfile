FROM golang:1.17.2-buster as builder
WORKDIR /app
ARG VERSION=dev
ARG GOFLAGS
COPY . /app
RUN make build

FROM debian:10.9-slim as final
RUN set -ex &&\
 apt-get update &&\
 DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates tzdata &&\
 rm -rf /var/lib/apt/lists/*
WORKDIR /data/inmemory-search
EXPOSE 8080 9000
ENTRYPOINT ["/data/benzinga-backend-challenge/benzinga-backend-challenge"]
COPY . /app
COPY --from=builder /app/build/* /data/benzinga-backend-challenge/

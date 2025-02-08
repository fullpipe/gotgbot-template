# syntax = docker.io/docker/dockerfile:experimental

# BUILDER
FROM golang as builder

WORKDIR /app

# BUILD
FROM builder as build

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/app

# release
FROM scratch as release

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/app /go/bin/app

ENTRYPOINT ["/go/bin/app"]

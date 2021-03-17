FROM golang:1.13 as build-env

ENV CONFIG_PATH $CONFIG_PATH

WORKDIR /app

COPY ./ .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /opt/cacher ./cmd/cacher/main.go

# Release
FROM alpine:latest

WORKDIR /root/

COPY --from=build-env /opt/cacher .
COPY --from=build-env /app/config ./config

CMD ["./cacher", "run", "$CONFIG_PATH"]

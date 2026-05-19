FROM golang:1.23-alpine AS build

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata git build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/kanvix-api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata libgcc

WORKDIR /app

COPY --from=build /out/kanvix-api /usr/local/bin/kanvix-api

EXPOSE 8080

ENV KANVIX_ADDR=:8080

CMD ["kanvix-api"]

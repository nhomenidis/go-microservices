FROM golang:1.21.2-alpine3.18 as base

RUN mkdir /app

COPY authentication /app/authentication
COPY common /app/common

WORKDIR /app/authentication

RUN CGO_ENABLED=0 GOOS=linux go build -o authApp ./cmd/api

RUN chmod +x /app/authentication/authApp

FROM alpine:3.18

RUN mkdir /app

COPY --from=base /app/authentication/authApp /app

CMD ["/app/authApp"]

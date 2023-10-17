FROM golang:1.21-alpine as base

RUN mkdir /app

COPY logger /app/logger
COPY common /app/common

WORKDIR /app/logger
RUN CGO_ENABLED=0 go build -o loggerApp ./cmd/api

RUN chmod +x /app/logger/loggerApp
# consider removing the build section
#  Final
From alpine:latest

RUN mkdir /app

COPY --from=base /app/logger/loggerApp /app

CMD ["/app/loggerApp"]
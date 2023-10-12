# base go image
# consider removing the build section

FROM golang:1.21-alpine as base

RUN mkdir /app

COPY broker /app/broker
COPY common /app/common

WORKDIR /app/broker
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/broker/brokerApp
# consider removing the build section
#  Final
From alpine:latest

RUN mkdir /app

COPY --from=base /app/broker/brokerApp /app

CMD ["/app/brokerApp"]

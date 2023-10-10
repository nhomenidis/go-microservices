# base go image
# consider removing the build section

FROM golang:1.21-alpine as base

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp
# consider removing the build section
#  Final
From alpine:latest

RUN mkdir /app

COPY --from=base /app/brokerApp /app

CMD ["/app/brokerApp"]

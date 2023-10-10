FROM golang:1.21.2-alpine3.18 as base

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o authApp ./cmd/api

RUN chmod +x /app/authApp

RUN apk --no-cache add file
RUN file /app/authApp > /file-output.txt

FROM alpine:3.18

# Copy the file-output.txt from the first stage
COPY --from=base /file-output.txt /file-output.txt

RUN mkdir /app

COPY --from=base /app/authApp /app

CMD ["/app/authApp"]

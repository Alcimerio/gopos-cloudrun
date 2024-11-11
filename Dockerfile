FROM golang:1.23.3 AS build

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun ./cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/cloudrun .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["./cloudrun"]
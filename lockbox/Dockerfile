FROM golang:1.14-alpine as builder

COPY go.mod /service/
COPY go.sum /service/

WORKDIR /service

RUN go mod download

COPY utils/*.go ./utils/
COPY rpc/*.go ./rpc/

COPY lockbox/*.go ./lockbox/

RUN GOOS=linux GOARCH=amd64 go build -o lockbox-service ./lockbox

FROM alpine

COPY --from=builder /service/lockbox-service /service/lockbox-service

WORKDIR /service

ENTRYPOINT ["./lockbox-service"]



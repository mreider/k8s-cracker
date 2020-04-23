FROM golang:1.14-alpine as builder

COPY go.mod /service/
COPY go.sum /service/

WORKDIR /service

RUN go mod download

COPY utils/*.go ./utils/
COPY rpc/*.go ./rpc/

COPY cracker/*.go ./cracker/

RUN GOOS=linux GOARCH=amd64 go build -o cracker-service ./cracker

FROM alpine

COPY --from=builder /service/cracker-service /service/cracker-service

WORKDIR /service

ENTRYPOINT ["./cracker-service"]



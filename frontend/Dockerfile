FROM golang:1.14-alpine as builder

COPY go.mod /service/
COPY go.sum /service/

WORKDIR /service

RUN go mod download

COPY utils/*.go ./utils/
COPY rpc/*.go ./rpc/

COPY frontend/*.go ./frontend/

RUN go build -o frontend-service ./frontend

FROM alpine

COPY --from=builder /service/frontend-service /service/frontend-service

WORKDIR /service

ENTRYPOINT ["./frontend-service"]



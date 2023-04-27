FROM golang:1.18-alpine

WORKDIR /client
COPY go.mod .
RUN go mod download -x

COPY . /client

RUN go build -o bin/server ./cmd/client

CMD ["./bin/server"]
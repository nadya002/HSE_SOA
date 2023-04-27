FROM golang:1.18-alpine

WORKDIR /server
COPY go.mod .
RUN go mod download -x

COPY . /server

RUN go build -o bin/server ./cmd/server

CMD ["./bin/server"]
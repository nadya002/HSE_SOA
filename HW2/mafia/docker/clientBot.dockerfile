FROM golang:1.18-alpine

WORKDIR /client
COPY go.mod .
RUN go mod download -x

COPY . /client

RUN go build -o bin/client ./cmd/client

#ARG format="json"
#ENV format
#ARG format
#ENV envFormat=$format

#EXPOSE 8080

#ENTRYPOINT ["./bin/server", "${format}"]
CMD ["sh", "-c", "./bin/client bot"]
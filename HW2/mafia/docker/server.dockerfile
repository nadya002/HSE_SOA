FROM golang:1.18-alpine

WORKDIR /server
COPY go.mod .
RUN go mod download -x

COPY . /server

RUN go build -o bin/server ./cmd/server



#ARG format="json"
#ENV format
#ARG format
#ENV envFormat=$format
ARG numbOfPl
ENV envNumbOfPl=$numbOfPl

EXPOSE 50005

#ENTRYPOINT ["./bin/server", "${format}"]
CMD ["sh", "-c", "./bin/server ${envNumbOfPl}"]

#ENTRYPOINT ["./bin/server ${envNumbOfPl}"]
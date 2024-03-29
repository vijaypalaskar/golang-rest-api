FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download && go mod verify

COPY . ./

RUN go build -o /golang-rest-api

EXPOSE 8000

CMD ["/golang-rest-api"]

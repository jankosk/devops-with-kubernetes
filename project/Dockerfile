FROM golang:alpine3.20

WORKDIR /app

COPY go.mod ./

RUN go mod download && go mod verify

COPY *.go .
COPY public /app/public

RUN go build -v -o /app/todo

CMD ["/app/todo"]
FROM golang:alpine3.20

WORKDIR /app

COPY go.mod ./

RUN go mod download && go mod verify

COPY *.go .

RUN go build -v -o /todo

CMD ["/todo"]
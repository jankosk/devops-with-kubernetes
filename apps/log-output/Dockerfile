FROM golang:alpine3.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY *.go .

RUN go build -v -o /log-output

CMD ["/log-output"]
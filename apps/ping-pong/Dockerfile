FROM golang:1.23-alpine AS build

WORKDIR /app/ping-pong

COPY apps/ping-pong/go.* .
COPY common/ /common/

RUN go mod download

COPY apps/ping-pong/*.go .

RUN CGO_ENABLED=0 go build -v -o /bin/ping-pong

FROM alpine

COPY --from=build /bin/ping-pong /
CMD ["/ping-pong"]
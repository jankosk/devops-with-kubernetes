FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY *.go .

RUN CGO_ENABLED=0 go build -v -o /bin/ticker

FROM gcr.io/distroless/static-debian12

COPY --from=build /bin/ticker /
CMD ["/ticker"]
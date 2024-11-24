FROM golang:1.23-alpine AS build

WORKDIR /app/log-output/ticker

COPY apps/log-output/ticker/ ./
COPY common/ /common/

RUN go mod download


RUN CGO_ENABLED=0 go build -v -o /bin/ticker

FROM alpine

COPY --from=build /bin/ticker /
CMD ["/ticker"]
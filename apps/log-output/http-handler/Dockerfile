FROM golang:1.23-alpine AS build

WORKDIR /app/log-output/http-handler

COPY apps/log-output/http-handler/ ./
COPY common/ /common/

RUN go mod download

RUN CGO_ENABLED=0 go build -v -o /bin/log-handler

FROM alpine

COPY --from=build /bin/log-handler /
CMD ["/log-handler"]
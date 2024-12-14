FROM golang:1.23-alpine AS build

WORKDIR /project/todo-api

COPY common/ /common/
COPY project/todo-api/go.* .

RUN go mod download

COPY project/todo-api/*.go .

RUN CGO_ENABLED=0 go build -v -o /bin/todo-api

FROM alpine

COPY --from=build /bin/todo-api /
CMD ["/todo-api"]
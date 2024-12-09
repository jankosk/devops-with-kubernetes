FROM golang:1.23-alpine AS build

WORKDIR /project/todo-app

COPY project/todo-app/go.mod .
COPY common/ /common/

RUN go mod download

COPY project/todo-app/*.go .
COPY project/todo-app/public/ ./public/

RUN CGO_ENABLED=0 go build -v -o /bin/todo

FROM alpine

COPY --from=build /bin/todo /
COPY --from=build project/todo-app/public/ /public/
CMD ["/todo"]
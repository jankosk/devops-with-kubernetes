FROM golang:1.23-alpine AS build

WORKDIR /project

COPY /project/ ./
COPY common/ /common/

RUN go mod download

RUN CGO_ENABLED=0 go build -v -o /bin/todo

FROM alpine

COPY --from=build /bin/todo /
COPY --from=build /project/public/ /public/
CMD ["/todo"]
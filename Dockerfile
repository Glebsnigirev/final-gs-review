FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN go build -o todo-server ./main.go

FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/todo-server .
COPY web ./web
ENV TODO_PORT=7540
ENV TODO_DBFILE=/db/todo.db
EXPOSE 7540
CMD ["./todo-server"]

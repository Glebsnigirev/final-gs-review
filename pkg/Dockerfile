
FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates

COPY todo-server /app/todo-server
COPY web /app/web

WORKDIR /app
EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/db/todo.db

CMD ["./todo-server"]

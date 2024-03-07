FROM golang:latest

RUN apt-get update && apt-get install -y sqlite3

RUN mkdir /forum

COPY . /forum

WORKDIR /forum

RUN go build -o main cmd/main.go

EXPOSE 8080

CMD ["/forum/main"]
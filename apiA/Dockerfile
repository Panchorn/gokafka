FROM golang:1.22-alpine AS builder

RUN mkdir /encryption
RUN mkdir /events
RUN mkdir /logs
COPY encryption/ /encryption
COPY events/ /events
COPY logs/ /logs

#ENV APP_PORT=8000
ENV KAFKA_SERVERS="kafka:9092"
ENV DB_HOST=mysql
ENV REDIS_HOST=redis

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8000
ENTRYPOINT ["./main"]

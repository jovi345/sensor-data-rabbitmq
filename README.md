# Sensor Data Processing with RabbitMQ

A program that sends sensor data (dummy) periodically every 1 minute.

## Installation

Install and run RabbitMQ with Docker for an easy process

```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

Install the dependency

```bash
go mod init <your-module-name>
go get github.com/rabbitmq/amqp091-go
go mod tidy
```

## Open RabbitMQ UI

Access from browser

```bash
localhost:15672/
```

Login with username & password: <b>guest</b>

## Run the program

```bash
go run publish/publish.go
```

```bash
go run subscribe/subscribe.go
```

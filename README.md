# kafka-stream 🧙🏼‍♂️

Service responsible for streaming Kafka messages.


## What it does?

This service reads all messages from the input topic and sends the messages to an SNS, using the lib [Goka](https://github.com/lovoo/goka) to stream messages.


## Goka?

Goka is a compact and powerful distributed stream processing library for Apache Kafka, written in Go. Goka aims to reduce the complexity of building highly scalable and highly available microservices.

Goka extends the concept of Kafka consumer groups by linking a state table to them and persisting them in Kafka itself. Goka provides logical patterns and a pluggable architecture.


## How to execute

- Execute the command `docker-compose up -d` to init Zookeeper, Kafka and Kowl.
- Execute the command `make create_topics` to create the input topic.
- Set the environment vars (env.example).
- Execute the command `make run` if you have godotenv installed or `go run main.go` to start the application.
- Send a message to Input topic.


## Kowl 

- Access the address `http://localhost:8080/` to access Kowl, a UI for Kafka, it will help you.

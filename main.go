package main

import (
	"log"

	"github.com/luansapelli/kafka-stream/clients/notification"
	"github.com/luansapelli/kafka-stream/controllers"
	"github.com/luansapelli/kafka-stream/environment"
	"github.com/luansapelli/kafka-stream/kafka"
	"github.com/luansapelli/kafka-stream/service"
)

func main() {
	env, err := environment.Load()
	if err != nil {
		log.Fatalf("error to load environment variables - %s", err)
	}

	config := kafka.Config(env.Application.Name).Sarama()
	producer := kafka.NewProducer(env.Kafka.Brokers, env.Kafka.Topic.Input, config)
	sns := notification.NewSNS(env.AWS.SNS.ARN)
	stream := service.InitStream(env, sns, config)

	doneChannel := make(chan bool)

	go controllers.HealthCheck()
	go stream.Run(doneChannel)

	producer.SendConfigMessages()

	doneChannel <- true
}

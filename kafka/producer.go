package kafka

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	kafkaBrokers []string
	inputTopic string
	config *sarama.Config
}

// NewProducer populate and returns a *Producer struct.
func NewProducer(kafkaBrokers []string, inputTopic string, config *sarama.Config) *Producer {
	return &Producer{
		kafkaBrokers: kafkaBrokers,
		inputTopic: inputTopic,
		config: config,
	}
}

// SendConfigMessages send one message to all partitions of input topic when the application starts.
// This method have 15 seconds of timeSleep to wait entire processor run
// and them, send the configuration messages.
func (producer *Producer) SendConfigMessages() {
	client := producer.createClient()
	syncProducer := producer.createSyncProducer(client)

	time.Sleep(15 * time.Second)

	producer.sendMessage(syncProducer, client)
}

func (producer *Producer) createClient() sarama.Client {
	client, err := sarama.NewClient(producer.kafkaBrokers, producer.config)
	if err != nil {
		log.Fatalf("error to create kafka client - %s", err)
	}
	return client
}

func (producer *Producer) createSyncProducer(kafkaClient sarama.Client) sarama.SyncProducer {
	syncProducer, err := sarama.NewSyncProducerFromClient(kafkaClient)
	if err != nil {
		log.Fatalf("error to create kafka sync producer - %s", err)
	}
	return syncProducer
}

func (producer *Producer) sendMessage(syncProducer sarama.SyncProducer, client sarama.Client) {
	partitions, err := client.Partitions(producer.inputTopic)
	if err != nil {
		log.Fatalf("error to get partitions from input topic - %s", err)
	}

	for i := range partitions {
		partition := partitions[i]

		msg := sarama.ProducerMessage{
			Topic:     producer.inputTopic,
			Value:     sarama.StringEncoder("this is a configuration message"),
			Key:       sarama.StringEncoder("configuration message"),
			Partition: partition,
		}

		p, _, err := syncProducer.SendMessage(&msg)
		if err != nil {
			log.Fatalf("error to send configuration message - %s", err)
		}

		log.Printf("sending configuration message to partition %v", p)
	}
}

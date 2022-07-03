package service

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"github.com/luansapelli/kafka-stream/clients/notification"
	"github.com/luansapelli/kafka-stream/environment"
)

type Stream struct {
	env          *environment.Environment
	notification notification.ServiceNotification
	saramaConfig *sarama.Config
	saramaLogger sarama.StdLogger
}

// InitStream populate and return a *Stream struct.
func InitStream(env *environment.Environment, notification notification.ServiceNotification, saramaConfig *sarama.Config) *Stream {
	sarama.Logger = log.New(os.Stdout, "", log.LstdFlags)

	return &Stream{
		env:          env,
		notification: notification,
		saramaConfig: saramaConfig,
		saramaLogger: sarama.Logger,
	}
}

// Run reads messages from kafka and send to a notification service.
func (stream *Stream) Run(doneChannel chan bool) {
	stream.saramaLogger.Println("starting streaming...")

	// callback is invoked for each message delivered from input topic.
	callback := func(context goka.Context, message interface{}) {
		contextKey := strings.Replace(context.Key(), `"`, "", -1)

		if contextKey == "configuration message" {
			stream.updateGroupTable(context)
			stream.saramaLogger.Printf("processing a configuration message - %s", message)
			return
		}

		rawMessage, err := stream.parseMessage(message)
		if err != nil {
			stream.updateGroupTable(context)
			stream.saramaLogger.Printf("error to parse message - %s", err)
			return
		}

		_, err = stream.notification.Write(rawMessage)
		if err != nil {
			stream.saramaLogger.Printf("error to send notification - %s", err)
		}

		stream.updateGroupTable(context)
	}

	processor := stream.createProcessor(callback)

	
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := processor.Run(ctx); err != nil {
			stream.saramaLogger.Printf("error to run processor - %s", err)
			os.Exit(1)
		}
		doneChannel <- true
	}()

	sig := make(chan os.Signal)
	go func() {
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	}()

	select {
	case <-sig:
	case <-done:
	}
	cancel()
	<-done
}

func (stream *Stream) createProcessor(callback func(context goka.Context, message interface{})) *goka.Processor {
	groupGraph := goka.DefineGroup(goka.Group(stream.env.Application.Name),
		goka.Input(goka.Stream(stream.env.Kafka.Topic.Input), new(codec.String), callback),
		goka.Persist(new(codec.Int64)))

	processor, err := goka.NewProcessor(stream.env.Kafka.Brokers,
		groupGraph,
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithConfig(stream.saramaConfig, goka.NewTopicManagerConfig())),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder))
	if err != nil {
		stream.saramaLogger.Printf("error to create processor - %s", err)
		os.Exit(1)
	}

	return processor
}

func (stream *Stream) parseMessage(message interface{}) ([]byte, error) {
	var messageMap map[string]interface{}

	rawMessage, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	strMessage, err := strconv.Unquote(string(rawMessage))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(strMessage), &messageMap)
	if err != nil {
		return nil, err
	}

	rawMessage, err = json.Marshal(messageMap)
	if err != nil {
		return nil, err
	}

	return rawMessage, nil
}

func (stream *Stream) updateGroupTable(context goka.Context) {
	var counter int64

	if val := context.Value(); val != nil {
		counter = val.(int64)
	}

	counter++
	context.SetValue(counter)
}

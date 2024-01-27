package client

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

const (
	KafkaServerAddress = "kafka:9092"
)
const (
	Topic = "movieWatchedByUser"
)

func SetupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress},
		config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}
func (c *Client) SendMovieWatchedMessage(movieId string, userId string) error {
	msg := &sarama.ProducerMessage{
		Topic: Topic,
		Key:   sarama.StringEncoder(userId),
		Value: sarama.StringEncoder(userId + ":" + movieId),
	}
	_, _, err := c.KafkaSamaraProducer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func ConsumeWatchedMovies(consumer sarama.ConsumerGroupHandler) error {
	consumerGroup, err := sarama.NewConsumerGroup([]string{KafkaServerAddress}, "watched-movies-suggestion", nil)
	if err != nil {
		return err
	}
	ctx := context.Background()
	for {
		err := consumerGroup.Consume(ctx, []string{Topic}, consumer)
		if err != nil {
			log.Panicf("Error from consumer: %v", err)
			return err
		}
	}
}

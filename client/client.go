package client

import "github.com/IBM/sarama"

type (
	Client struct {
		KafkaSamaraProducer sarama.SyncProducer
	}
)

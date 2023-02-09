package kafka

import (
	"github.com/IBM/sarama"
)

type KafkaClient struct {
	kAdmin  sarama.ClusterAdmin
	kClient sarama.Client
}

func NewKafkaClient() (*KafkaClient, error) {
	kClient, err := sarama.NewClient([]string{"localhost:9093"}, nil)
	if err != nil {
		return nil, err
	}

	kAdmin, err := sarama.NewClusterAdminFromClient(kClient)
	if err != nil {
		return nil, err
	}

	return &KafkaClient{
		kAdmin:  kAdmin,
		kClient: kClient,
	}, nil
}

func (k *KafkaClient) Close() error {

	if err := k.kAdmin.Close(); err != nil {
		return err
	}

	if err := k.kClient.Close(); err != nil {
		return err
	}

	return nil
}

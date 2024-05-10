package kafka

import (
	"strconv"

	"github.com/IBM/sarama"
)

func (k *Service) Version() string {
	metadata, _ := k.kClient.Brokers()[0].GetMetadata(&sarama.MetadataRequest{})
	return strconv.Itoa(int(metadata.Version))

}

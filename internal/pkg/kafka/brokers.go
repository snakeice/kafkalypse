package kafka

import (
	"strconv"

	"github.com/IBM/sarama"
)

type Broker struct {
	ID     int32
	Host   string
	Leader bool
	// TODO: add more fields?
}

func (k *Service) Version() string {
	metadata, _ := k.kClient.Brokers()[0].GetMetadata(&sarama.MetadataRequest{})
	return strconv.Itoa(int(metadata.Version))

}
func (k *Service) BrokersStr() []string {
	brokers := make([]string, 0)

	for _, broker := range k.kClient.Brokers() {
		brokers = append(brokers, broker.Addr())
	}

	return brokers
}

func (k *Service) Brokers() []*Broker {
	brokers := make([]*Broker, 0)

	for _, broker := range k.kClient.Brokers() {
		brokers = append(brokers, &Broker{
			ID:     broker.ID(),
			Host:   broker.Addr(),
			Leader: false,
		})
	}

	return brokers

}

package kafka

import (
	"github.com/IBM/sarama"
)

type TopicOverview struct {
	Name        string
	Partitions  int32
	Replication int16
	IsInternal  bool
}

type TopicDetail struct {
	*TopicOverview
	Config map[string]string
}

func (k *KafkaClient) ListTopics() ([]TopicOverview, error) {
	topics, err := k.kAdmin.ListTopics()
	if err != nil {
		return nil, err
	}

	var topicsOverview []TopicOverview

	for topic, detail := range topics {
		topicsOverview = append(topicsOverview, TopicOverview{
			Name:        topic,
			Partitions:  detail.NumPartitions,
			Replication: detail.ReplicationFactor,
			IsInternal:  false, // TODO: check if topic is internal
		})
	}

	return topicsOverview, nil
}

func (k *KafkaClient) TopicDetail(topic *TopicOverview) (*TopicDetail, error) {
	var topicDetail TopicDetail

	topicDetail.TopicOverview = topic

	config, err := k.kAdmin.DescribeConfig(sarama.ConfigResource{
		Type: sarama.TopicResource,
		Name: topic.Name,
	})
	if err != nil {
		return nil, err
	}

	topicDetail.Config = make(map[string]string, len(config))
	for _, entry := range config {
		topicDetail.Config[entry.Name] = entry.Value
	}

	return &topicDetail, nil

}

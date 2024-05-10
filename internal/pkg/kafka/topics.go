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

func (k *Service) ListTopics() ([]TopicOverview, error) {
	topics, err := k.kAdmin.ListTopics()
	if err != nil {
		return nil, err
	}

	var topicsOverview []TopicOverview

	for topic, detail := range topics {
		descriptions, err := k.kAdmin.DescribeTopics([]string{topic})
		if err != nil {
			return nil, err
		}

		desc := descriptions[0]

		topicsOverview = append(topicsOverview, TopicOverview{
			Name:        topic,
			Partitions:  detail.NumPartitions,
			Replication: detail.ReplicationFactor,
			IsInternal:  desc.IsInternal,
		})
	}

	return topicsOverview, nil
}

func (k *Service) TopicDetail(topic *TopicOverview) (*TopicDetail, error) {
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

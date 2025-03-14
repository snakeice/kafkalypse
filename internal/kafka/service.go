package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/config"
	"github.com/snakeice/kafkalypse/internal/tools"
)

type Service struct {
	kAdmin  sarama.ClusterAdmin
	kClient sarama.Client
	ctx     context.Context
	version string
}

func NewKafkaClient(ctxConfig *config.KafkaContext) (*Service, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0 // Initial version for connection
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true

	kClient, err := sarama.NewClient(strings.Split(ctxConfig.BootstrapServers, ","), config)
	if err != nil {
		return nil, err
	}

	// Get actual broker version
	broker := kClient.Brokers()[0]

	if ok, err := broker.Connected(); !ok {
		if err != nil {
			return nil, err
		}

		if err := broker.Open(config); err != nil {
			return nil, err
		}
	}

	// Reconnect with correct version
	kClient.Close()
	kClient, err = sarama.NewClient(strings.Split(ctxConfig.BootstrapServers, ","), config)
	if err != nil {
		return nil, err
	}

	kAdmin, err := sarama.NewClusterAdminFromClient(kClient)
	if err != nil {
		return nil, err
	}

	return &Service{
		kAdmin:  kAdmin,
		kClient: kClient,
		ctx:     context.Background(),
		version: ":TODO:",
	}, nil
}

func (k *Service) Close() error {
	if err := k.kAdmin.Close(); err != nil {
		return err
	}

	if err := k.kClient.Close(); err != nil {
		return err
	}

	return nil
}

func (k *Service) BrokersStr() string {
	brokers := k.kClient.Brokers()
	var strs []string
	for _, b := range brokers {
		strs = append(strs, b.Addr())
	}
	return strings.Join(strs, ",")
}

func (k *Service) Version() string {
	return k.version
}

// Event Handlers
func (k *Service) handleTopicListEvent() tea.Msg {
	topics, err := k.kAdmin.ListTopics()
	if err != nil {
		return TopicListEventRes{Error: err}
	}

	var topicsOverview []TopicOverview
	topicsNames := make([]string, 0, len(topics))
	for topicName := range topics {
		topicsNames = append(topicsNames, topicName)
	}

	details, err := k.kAdmin.DescribeTopics(topicsNames)
	if err != nil {
		return TopicListEventRes{Error: err}
	}

	for topicName, topic := range topics {
		for _, detail := range details {
			if detail.Name == topicName {
				topicsOverview = append(topicsOverview, TopicOverview{
					Name:        detail.Name,
					Partitions:  int32(len(detail.Partitions)),
					Replication: topic.ReplicationFactor,
					IsInternal:  detail.IsInternal,
				})
				break
			}
		}
	}

	return TopicListEventRes{Topics: topicsOverview}
}

func (k *Service) handleTopicDetailEvent(event TopicDetailEventReq) tea.Msg {
	config, err := k.kAdmin.DescribeConfig(sarama.ConfigResource{
		Type: sarama.TopicResource,
		Name: event.Topic,
	})
	if err != nil {
		return TopicDetailEventRes{Error: err}
	}

	topicDetail := &TopicDetail{
		TopicOverview: TopicOverview{
			Name: event.Topic,
		},
		Config: make(map[string]string, len(config)),
	}

	for _, entry := range config {
		topicDetail.Config[entry.Name] = entry.Value
	}

	return TopicDetailEventRes{Topic: topicDetail}
}

func (k *Service) handleTopicCreateEvent(event TopicCreateEventReq) tea.Msg {
	err := k.kAdmin.CreateTopic(event.Topic, &sarama.TopicDetail{
		NumPartitions:     event.Partitions,
		ReplicationFactor: event.Replication,
	}, false)
	return TopicCreateEventRes{Topic: event.Topic, Error: err}
}

func (k *Service) handleTopicDeleteEvent(event TopicDeleteEventReq) tea.Msg {
	err := k.kAdmin.DeleteTopic(event.Topic)
	return TopicDeleteEventRes{Topic: event.Topic, Error: err}
}

func (k *Service) handleTopicConfigUpdateEvent(event TopicConfigUpdateEventReq) tea.Msg {
	configMap := make(map[string]*string)
	for key, value := range event.Config {
		val := value // Create a new variable to get a unique address
		configMap[key] = &val
	}

	err := k.kAdmin.AlterConfig(sarama.TopicResource, event.Topic, configMap, false)
	if err != nil {
		return TopicConfigUpdateEventRes{Error: err}
	}

	return TopicConfigUpdateEventRes{Topic: event.Topic, Config: event.Config}
}

func (k *Service) handleConsumerGroupListEvent() tea.Msg {
	groups, err := k.kAdmin.ListConsumerGroups()
	if err != nil {
		return ConsumerGroupListEvent{Error: err}
	}

	var consumerGroups []ConsumerGroup
	for groupID := range groups {
		consumerGroups = append(consumerGroups, ConsumerGroup{
			GroupID: groupID,
			State:   groups[groupID],
		})
	}

	return ConsumerGroupListEvent{Groups: consumerGroups}
}

func (k *Service) handleMessageListEvent(event MessageListEvent) tea.Msg {
	consumer, err := sarama.NewConsumerFromClient(k.kClient)
	if err != nil {
		return MessageListEvent{Error: err}
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(event.Topic)
	if err != nil {
		return MessageListEvent{Error: err}
	}

	var messages []Message
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(event.Topic, partition, sarama.OffsetNewest)
		if err != nil {
			continue
		}

		select {
		case msg := <-pc.Messages():
			messages = append(messages, Message{
				Key:       msg.Key,
				Value:     msg.Value,
				Offset:    msg.Offset,
				Timestamp: msg.Timestamp.Unix(),
			})
		case <-time.After(100 * time.Millisecond):
		}
		pc.Close()
	}

	return MessageListEvent{
		Topic:    event.Topic,
		Messages: messages,
	}
}

func (k *Service) handleMessageProduceEvent(event MessageProduceEvent) tea.Msg {
	producer, err := sarama.NewSyncProducerFromClient(k.kClient)
	if err != nil {
		return MessageProduceEvent{Error: err}
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: event.Topic,
		Value: sarama.StringEncoder("test message"),
	}

	_, _, err = producer.SendMessage(msg)
	return MessageProduceEvent{Topic: event.Topic, Error: err}
}

func (k *Service) handleGetConnMsg() tea.Msg {
	// Check if we have a valid connection
	if k.kClient == nil || k.kAdmin == nil {
		return KafkaConnectionRes{
			KafkaService: nil,
			Err:          fmt.Errorf("no active connection"),
		}
	}

	// Check if we can connect to any broker
	brokers := k.kClient.Brokers()
	if len(brokers) == 0 {
		return KafkaConnectionRes{
			KafkaService: nil,
			Err:          fmt.Errorf("no brokers available"),
		}
	}

	// Try to connect to the first broker
	broker := brokers[0]
	if ok, err := broker.Connected(); !ok {
		if err != nil {
			return KafkaConnectionRes{
				KafkaService: nil,
				Err:          err,
			}
		}
	}

	return KafkaConnectionRes{
		KafkaService: k,
		Err:          nil,
	}
}

func (k *Service) Update(msg tea.Msg) tea.Cmd {
	var result tea.Msg
	switch msg := msg.(type) {
	case TopicListEventReq:
		result = k.handleTopicListEvent()
	case TopicDetailEventReq:
		result = k.handleTopicDetailEvent(msg)
	case TopicCreateEventReq:
		result = k.handleTopicCreateEvent(msg)
	case TopicDeleteEventReq:
		result = k.handleTopicDeleteEvent(msg)
	case TopicConfigUpdateEventReq:
		result = k.handleTopicConfigUpdateEvent(msg)
	case ConsumerGroupListEvent:
		result = k.handleConsumerGroupListEvent()
	case MessageListEvent:
		result = k.handleMessageListEvent(msg)
	case MessageProduceEvent:
		result = k.handleMessageProduceEvent(msg)
	case KafkaConnectionReq:
		result = k.handleGetConnMsg()
	default:
		return nil
	}

	return tools.WrapCmd(result)
}

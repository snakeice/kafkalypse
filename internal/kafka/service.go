package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/config"
	"github.com/snakeice/kafkalypse/internal/tools"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Service struct {
	client  *kgo.Client
	admin   *kadm.Client
	ctx     context.Context
	version string
}

func NewKafkaClient(ctxConfig *config.KafkaContext) (*Service, error) {
	bootstrapServers := strings.Split(ctxConfig.BootstrapServers, ",")
	ctx := context.Background()

	// Configure the Kafka client
	opts := []kgo.Opt{
		kgo.SeedBrokers(bootstrapServers...),
		kgo.ClientID("kafkalypse"),
		kgo.DialTimeout(10 * time.Second),
	}

	// Add security options if configured
	if ctxConfig.SecurityProtocol != "" {
		switch ctxConfig.SecurityProtocol {
		case "SASL_PLAINTEXT":
			opts = append(opts,
				kgo.SASL(ctxConfig.Sasl.AsMechanism()),
			)
		case "SASL_SSL":
			tlsDialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}

			opts = append(opts,
				kgo.SASL(ctxConfig.Sasl.AsMechanism()),
				kgo.Dialer(tlsDialer.DialContext),
			)
		}
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	admin := kadm.NewClient(client)

	apiVersions, err := admin.ApiVersions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get API versions: %v", err)
	}

	ver := apiVersions.Sorted()[0]

	// Get broker version
	var version string
	if ver.Err != nil {
		version = "Unknown"
	} else {
		version = ver.VersionGuess()
	}

	return &Service{
		client:  client,
		admin:   admin,
		ctx:     ctx,
		version: version,
	}, nil
}

func (k *Service) Close() error {
	k.client.Close()
	return nil
}

func (k *Service) Brokers() []string {
	// Get brokers from metadata
	req := kmsg.NewMetadataRequest()
	resp, err := req.RequestWith(k.ctx, k.client)
	if err != nil {
		return nil
	}

	var brokers []string
	for _, broker := range resp.Brokers {
		brokers = append(brokers, fmt.Sprintf("%s:%d", broker.Host, broker.Port))
	}

	return brokers
}

func (k *Service) Version() string {
	return k.version
}

// Event Handlers
func (k *Service) handleTopicListEvent() tea.Msg {
	// List topics using kadm
	topics, err := k.admin.ListTopics(k.ctx)
	if err != nil {
		return TopicListEventRes{Error: err}
	}

	var topicsOverview []TopicOverview
	for name, topic := range topics {

		topicsOverview = append(topicsOverview, TopicOverview{
			Name:        name,
			Partitions:  len(topic.Partitions),
			Replication: topic.Partitions.NumReplicas(),
			IsInternal:  topic.IsInternal,
		})
	}

	return TopicListEventRes{Topics: topicsOverview}
}

func (k *Service) handleTopicDetailEvent(event TopicDetailEventReq) tea.Msg {
	configs, err := k.admin.DescribeTopicConfigs(k.ctx, event.Topic)
	if err != nil {
		return TopicDetailEventRes{Error: err}
	}

	if len(configs) == 0 {
		return TopicDetailEventRes{Error: fmt.Errorf("topic not found")}
	}

	config := configs[0]
	topicDetail := &TopicDetail{
		TopicOverview: TopicOverview{
			Name: event.Topic,
		},
		Config: make(map[string]string, len(config.Configs)),
	}

	for _, entry := range config.Configs {
		if entry.Value != nil {
			topicDetail.Config[entry.Key] = *entry.Value
		}
	}

	return TopicDetailEventRes{Topic: topicDetail}
}

func (k *Service) handleTopicCreateEvent(event TopicCreateEventReq) tea.Msg {
	_, err := k.admin.CreateTopic(
		k.ctx,
		event.Partitions,
		event.Replication,
		nil,
		event.Topic,
	)

	return TopicCreateEventRes{Topic: event.Topic, Error: err}
}

func (k *Service) handleTopicDeleteEvent(event TopicDeleteEventReq) tea.Msg {
	res, err := k.admin.DeleteTopic(k.ctx, event.Topic)
	if res.Err != nil {
		err = res.Err
	}

	return TopicDeleteEventRes{Topic: event.Topic, Error: err}
}

func (k *Service) handleTopicConfigUpdateEvent(event TopicConfigUpdateEventReq) tea.Msg {
	var updates []kadm.AlterConfig

	for name, value := range event.Config {
		updates = append(updates, kadm.AlterConfig{
			Name:  name,
			Op:    kadm.SetConfig, // TODO: handle other operations
			Value: &value,
		})
	}

	_, err := k.admin.AlterTopicConfigs(k.ctx, updates, event.Topic) //TODO: handle response
	if err != nil {
		return TopicConfigUpdateEventRes{Error: err}
	}

	return TopicConfigUpdateEventRes{Topic: event.Topic, Config: event.Config}
}

func (k *Service) handleConsumerGroupListEvent() tea.Msg {
	groups, err := k.admin.ListGroups(k.ctx)
	if err != nil {
		return ConsumerGroupListEvent{Error: err}
	}

	var consumerGroups []ConsumerGroup
	for id, group := range groups {
		consumerGroups = append(consumerGroups, ConsumerGroup{
			GroupID: id,
			State:   string(group.State),
		})
	}

	return ConsumerGroupListEvent{Groups: consumerGroups}
}

func (k *Service) handleMessageListEvent(event MessageListEvent) tea.Msg {

	// Set timeout for context
	ctx, cancel := context.WithTimeout(k.ctx, 500*time.Millisecond)
	defer cancel()

	var messages []Message

	// Poll for messages with a timeout
	fetches := k.client.PollFetches(ctx)
	errs := fetches.Errors()
	if len(errs) > 0 {
		return MessageListEvent{Error: errs[0].Err}
	}

	fetches.EachPartition(func(p kgo.FetchTopicPartition) {
		p.EachRecord(func(record *kgo.Record) {
			messages = append(messages, Message{
				Key:       record.Key,
				Value:     record.Value,
				Offset:    record.Offset,
				Timestamp: record.Timestamp.Unix(),
			})
		})
	})

	return MessageListEvent{
		Topic:    event.Topic,
		Messages: messages,
	}
}

func (k *Service) handleMessageProduceEvent(event MessageProduceEvent) tea.Msg {
	record := &kgo.Record{
		Topic: event.Topic,
		Value: []byte("test message"),
	}

	err := k.client.ProduceSync(k.ctx, record).FirstErr()
	return MessageProduceEvent{Topic: event.Topic, Error: err}
}

func (k *Service) handleGetConnMsg() tea.Msg {
	// Check if we have a valid connection
	if k.client == nil || k.admin == nil {
		return KafkaConnectionRes{
			KafkaService: nil,
			Err:          fmt.Errorf("no active connection"),
		}
	}

	// Try a simple metadata request to check connectivity
	err := k.client.Ping(k.ctx)
	if err != nil {
		return KafkaConnectionRes{
			KafkaService: nil,
			Err:          err,
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

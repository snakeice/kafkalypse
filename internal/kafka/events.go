package kafka

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/tools"
)

// Types
type Message struct {
	Key       []byte
	Value     []byte
	Offset    int64
	Timestamp int64
}

type TopicOverview struct {
	Name        string
	Partitions  int
	Replication int
	IsInternal  bool
}

type TopicDetail struct {
	TopicOverview
	Config map[string]string
}

type ConsumerGroup struct {
	GroupID string
	State   string
}

type ConsumerGroupDetail struct {
	GroupID string
	State   string
	Members []string
	Topics  []string
	Lag     map[string]int64
}

// Events
type Event interface {
	tea.Msg
}

type TickMsg time.Time

// Topic Events
type TopicListEventReq struct{}
type TopicListEventRes struct {
	Topics []TopicOverview
	Error  error
}

type TopicDetailEventReq struct {
	Topic string
}
type TopicDetailEventRes struct {
	Topic *TopicDetail
	Error error
}

type TopicCreateEventReq struct {
	Topic       string
	Partitions  int32
	Replication int16
}
type TopicCreateEventRes struct {
	Topic string
	Error error
}

type TopicDeleteEventReq struct {
	Topic string
}
type TopicDeleteEventRes struct {
	Topic string
	Error error
}

type TopicConfigUpdateEventReq struct {
	Topic  string
	Config map[string]string
}
type TopicConfigUpdateEventRes struct {
	Topic  string
	Config map[string]string
	Error  error
}

// Consumer Group Events
type ConsumerGroupListEvent struct {
	Groups []ConsumerGroup
	Error  error
}

type ConsumerGroupDetailEvent struct {
	Group *ConsumerGroupDetail
	Error error
}

// Message Events
type MessageListEvent struct {
	Topic    string
	Messages []Message
	Error    error
}

type MessageProduceEvent struct {
	Topic string
	Error error
}

func ListTopics() tea.Cmd {
	return tools.WrapCmd(TopicListEventReq{})
}

func GetTopicDetail(topic string) tea.Cmd {
	return tools.WrapCmd(TopicDetailEventReq{
		Topic: topic,
	})
}

func CreateTopic(name string, partitions int32, replicationFactor int16) tea.Cmd {
	return tools.WrapCmd(TopicCreateEventReq{
		Topic:       name,
		Partitions:  partitions,
		Replication: replicationFactor,
	})
}

func DeleteTopic(name string) tea.Cmd {
	return tools.WrapCmd(TopicDeleteEventReq{
		Topic: name,
	})
}

func ListConsumerGroups() tea.Cmd {
	return tools.WrapCmd(ConsumerGroupListEvent{
		Groups: []ConsumerGroup{},
	})
}

func GetConsumerGroupDetail(group string) tea.Cmd {
	return tools.WrapCmd(ConsumerGroupDetailEvent{
		Group: &ConsumerGroupDetail{
			GroupID: group,
		},
	})
}

func ListMessages(topic string, limit int) tea.Cmd {
	return tools.WrapCmd(MessageListEvent{
		Topic:    topic,
		Messages: []Message{},
	})
}

func ProduceMessage(topic string, message []byte) tea.Cmd {
	return tools.WrapCmd(MessageProduceEvent{
		Topic: topic,
	})
}

func Tick() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func UpdateTopicConfig(topic string, config map[string]string) tea.Cmd {
	return tools.WrapCmd(TopicConfigUpdateEventReq{
		Topic:  topic,
		Config: config,
	})
}

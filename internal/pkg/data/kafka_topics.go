package data

import (
	"fmt"
	"sort"
	"time"

	"github.com/snakeice/kafkalypse/internal/pkg/kafka"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/table"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/shortcut"
)

type KafkaTopics struct {
	conn   *kafka.KafkaClient
	active bool

	topics []kafka.TopicOverview
}

func NewKafkaTopics(conn *kafka.KafkaClient) *KafkaTopics {
	k := &KafkaTopics{
		conn:   conn,
		active: true,
	}

	go k.refreshLoop()
	return k
}

func (k *KafkaTopics) Len() int {

	return len(k.topics)
}

func (k *KafkaTopics) At(i int) []string {
	return []string{
		k.topics[i].Name,
		fmt.Sprint(k.topics[i].Partitions),
		fmt.Sprint(k.topics[i].Replication),
	}
}

func (k *KafkaTopics) Shortcuts() []shortcut.Action {
	return []shortcut.Action{
		// {Action: Insert, Shortcuts: []string{"i"}, Description: "insert"},
		// {Action: Delete, Shortcuts: []string{"d"}, Description: "describe"},
	}
}

func (k *KafkaTopics) Cols() []table.ColHead {
	return []table.ColHead{
		{Name: "Name", Perc: .6},
		{Name: "Partitions", Perc: .2},
		{Name: "Replication", Perc: .2},
	}
}

func (k *KafkaTopics) Refresh() error {
	topics, err := k.conn.ListTopics()
	if err != nil {
		return err
	}

	k.topics = topics

	sort.Slice(k.topics, func(i, j int) bool {
		return k.topics[i].Name < k.topics[j].Name
	})

	return nil
}

func (k *KafkaTopics) Active() bool     { return k.active }
func (k *KafkaTopics) SetActive(a bool) { k.active = a }

func (k *KafkaTopics) Close() error {
	return nil
}

func (k *KafkaTopics) Name() string {
	return "Topics"
}

func (k *KafkaTopics) refreshLoop() {
	for {
		if k.active {
			err := k.Refresh()
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

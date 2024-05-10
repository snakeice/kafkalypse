package topics

import (
	"sort"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/pkg/kafka"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/table"
)

type Topics struct {
	tbl      table.Model
	kafkaSvc *kafka.Service

	topicList []kafka.TopicOverview
}

func NewTopics() *Topics {
	t := &Topics{}
	t.tbl = table.NewTable(t)
	return t
}

func (t Topics) Init() tea.Cmd {
	return tea.Batch(
		t.tbl.Init(),
		kafka.GetConn(),
	)
}
func (t Topics) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case kafka.KafkaConnectionMsg:
		if msg.Err != nil {
			return t, nil
		}

		t.kafkaSvc = msg.Svc

		topicList, err := t.kafkaSvc.ListTopics()
		if err != nil {
			return t, nil
		}

		t.topicList = topicList
		sort.Slice(t.topicList, func(i, j int) bool {
			return t.topicList[i].Name < t.topicList[j].Name
		})

		t.tbl.Datasource = &t

		cmds = append(cmds, table.UpdateCmd())
	}

	tblModel, cmd := t.tbl.Update(msg)
	t.tbl = tblModel.(table.Model)
	cmds = append(cmds, cmd)
	return t, tea.Batch(cmds...)
}

func (t Topics) View() string {
	return t.tbl.View()
}

func (t *Topics) Len() int {
	return len(t.topicList)
}

func boolToStr(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func (t *Topics) At(idx int) []string {
	return []string{
		t.topicList[idx].Name,
		strconv.Itoa(int(t.topicList[idx].Partitions)),
		strconv.Itoa(int(t.topicList[idx].Replication)),
		boolToStr(t.topicList[idx].IsInternal),
	}
}
func (t *Topics) Cols() []table.ColHead {
	return []table.ColHead{
		{Name: "Topic", Perc: 0.75},
		{Name: "Partitions", Perc: 0.10},
		{Name: "Replicas", Perc: 0.10},
		{Name: "Internal", Perc: 0.05},
	}
}

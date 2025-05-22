package topics

import (
	"sort"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/kafka"
	"github.com/snakeice/kafkalypse/internal/tui/components/popup"
	"github.com/snakeice/kafkalypse/internal/tui/components/shortcuts"
	"github.com/snakeice/kafkalypse/internal/tui/components/table"
)

type Topics struct {
	tbl         table.Model
	kafkaSvc    *kafka.Service
	topicList   []kafka.TopicOverview
	error       error
	lastRefresh time.Time
	details     *TopicDetails
	shortcuts   *shortcuts.Shortcuts
	createPopup popup.Model
}

func NewTopics() *Topics {
	t := &Topics{
		lastRefresh: time.Now(),
		shortcuts:   shortcuts.New(),
	}
	t.tbl = table.NewTable(t)
	t.tbl.Focus() // Focus the table by default

	// Initialize create popup
	t.createPopup = popup.New("Create New Topic", []popup.Field{
		{Name: "Name", Value: "", Placeholder: "Enter topic name"},
		{Name: "Partitions", Value: "1", Placeholder: "Number of partitions"},
		{Name: "Replication Factor", Value: "1", Placeholder: "Replication factor"},
	})
	t.createPopup.Hide()

	// Add shortcuts for main view
	t.shortcuts.Add(shortcuts.Action{
		Key:         "n",
		Description: "new topic",
		Handler: func() tea.Cmd {
			t.createPopup.Show()
			return nil
		},
	})

	t.shortcuts.Add(shortcuts.Action{
		Key:         "d",
		Description: "delete topic",
		Handler: func() tea.Cmd {
			if t.tbl.Cursor() >= 0 && t.tbl.Cursor() < len(t.topicList) {
				topic := t.topicList[t.tbl.Cursor()]
				return kafka.DeleteTopic(topic.Name)
			}
			return nil
		},
	})

	t.shortcuts.Add(shortcuts.Action{
		Key:         "r",
		Description: "refresh",
		Handler:     kafka.ListTopics,
	})

	t.shortcuts.Add(shortcuts.Action{
		Key: "enter",
		Handler: func() tea.Cmd {
			if t.tbl.Cursor() >= 0 && t.tbl.Cursor() < len(t.topicList) {
				topic := t.topicList[t.tbl.Cursor()]
				t.details = NewTopicDetails(topic, t.kafkaSvc)
				return t.details.Init()
			}
			return nil
		},
	})

	return t
}

func (t Topics) Init() tea.Cmd {
	return tea.Batch(
		t.tbl.Init(),
		kafka.GetConn(),
		kafka.ListTopics(),
		kafka.Tick(),
	)
}

func (t Topics) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	if t.details != nil {
		model, cmd := t.details.Update(msg)

		//nolint:all
		if details, ok := model.(*TopicDetails); ok {
			t.details = details
		}

		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "q" && !t.details.editing {
			t.details = nil
			return t, nil
		}
		return t, tea.Batch(cmds...)
	}

	if t.createPopup.IsVisible() {
		model, cmd := t.createPopup.Update(msg)
		if updatedPopup, ok := model.(popup.Model); ok {
			t.createPopup = updatedPopup
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return t, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case popup.SubmitMsg:
		if name, ok := msg.Values["Name"]; ok && name != "" {
			partitions, _ := strconv.Atoi(msg.Values["Partitions"])
			if partitions < 1 {
				partitions = 1
			}
			replication, _ := strconv.Atoi(msg.Values["Replication Factor"])
			if replication < 1 {
				replication = 1
			}
			cmds = append(cmds, kafka.CreateTopic(name, int32(partitions), int16(replication)))
		}

	// case kafka.TickMsg:
	// 	// Auto-refresh topics
	// 	t.lastRefresh = time.Time(msg)
	// 	cmds = append(cmds, kafka.ListTopics())
	// 	// Schedule next tick
	// 	cmds = append(cmds, kafka.Tick())

	case kafka.KafkaConnectionRes:
		if msg.Err != nil {
			t.error = msg.Err
			return t, nil
		}

		t.kafkaSvc = msg.KafkaService
		cmds = append(cmds, kafka.ListTopics())

	case kafka.TopicListEventRes:
		if msg.Error != nil {
			t.error = msg.Error
			return t, nil
		}

		t.lastRefresh = time.Now()
		t.topicList = msg.Topics
		sort.Slice(t.topicList, func(i, j int) bool {
			return t.topicList[i].Name < t.topicList[j].Name
		})

		t.tbl.Datasource = &t
		cmds = append(cmds, table.UpdateCmd())

	case kafka.TopicCreateEventRes:
		if msg.Error != nil {
			t.error = msg.Error
			return t, nil
		}
		cmds = append(cmds, kafka.ListTopics())

	case kafka.TopicDeleteEventRes:
		if msg.Error != nil {
			t.error = msg.Error
			return t, nil
		}
		cmds = append(cmds, kafka.ListTopics())

	case tea.KeyMsg:
		// First try to handle navigation keys in the table
		tblModel, cmd := t.tbl.Update(msg)
		t.tbl = tblModel.(table.Model)
		if cmd != nil {
			cmds = append(cmds, cmd)
			return t, tea.Batch(cmds...)
		}

		// If table didn't handle the key and we're not editing, try shortcuts
		if t.details == nil || !t.details.editing {
			if cmd := t.shortcuts.Handle(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	tbl, cmd := t.tbl.Update(msg)
	t.tbl = tbl.(table.Model)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return t, tea.Batch(cmds...)
}

func (t Topics) View() string {
	if t.details != nil {
		return t.details.View() + "\n\nPress 'q' to return to topic list"
	}

	// Show popup if visible
	if popup := t.createPopup.View(); popup != "" {
		return popup
	}

	if t.error != nil {
		return "Error: " + t.error.Error()
	}
	if len(t.topicList) == 0 {
		return "No topics found. " + "\nLast refresh: " + t.lastRefresh.Format("15:04:05")
	}
	return t.tbl.View() + "\nLast refresh: " + t.lastRefresh.Format("15:04:05") + "\n" + t.shortcuts.GetHelpText()
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

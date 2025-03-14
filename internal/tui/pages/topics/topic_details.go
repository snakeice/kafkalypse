package topics

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/kafka"
)

type TopicDetails struct {
	topic     *kafka.TopicDetail
	kafkaSvc  *kafka.Service
	error     error
	editing   bool
	jsonInput string
	cursor    int
}

func NewTopicDetails(topic kafka.TopicOverview, kafkaSvc *kafka.Service) *TopicDetails {
	return &TopicDetails{
		topic: &kafka.TopicDetail{
			TopicOverview: topic,
		},
		kafkaSvc: kafkaSvc,
	}
}

func (t *TopicDetails) Init() tea.Cmd {
	return kafka.GetTopicDetail(t.topic.Name)
}

func (t *TopicDetails) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case kafka.TopicDetailEventRes:
		if msg.Error != nil {
			t.error = msg.Error
			return t, nil
		}
		t.topic = msg.Topic
		if t.editing {
			// Convert config to JSON for editing
			configJSON, err := json.MarshalIndent(t.topic.Config, "", "  ")
			if err != nil {
				t.error = err
				return t, nil
			}
			t.jsonInput = string(configJSON)
		}
		return t, nil

	case kafka.TopicConfigUpdateEventRes:
		if msg.Error != nil {
			t.error = msg.Error
			return t, nil
		}
		// Refresh topic details after successful update
		return t, kafka.GetTopicDetail(t.topic.Name)

	case tea.KeyMsg:
		// If we're editing, handle all key events
		if t.editing {
			switch msg.String() {
			case "esc":
				t.editing = false
				return t, nil
			case "enter":
				var newConfig map[string]string
				if err := json.Unmarshal([]byte(t.jsonInput), &newConfig); err != nil {
					t.error = err
					return t, nil
				}
				t.topic.Config = newConfig
				t.editing = false
				return t, kafka.UpdateTopicConfig(t.topic.Name, newConfig)
			case "backspace":
				if len(t.jsonInput) > 0 && t.cursor > 0 {
					t.jsonInput = t.jsonInput[:t.cursor-1] + t.jsonInput[t.cursor:]
					t.cursor--
				}
			case "left":
				if t.cursor > 0 {
					t.cursor--
				}
			case "right":
				if t.cursor < len(t.jsonInput) {
					t.cursor++
				}
			default:
				if len(msg.String()) == 1 {
					t.jsonInput = t.jsonInput[:t.cursor] + msg.String() + t.jsonInput[t.cursor:]
					t.cursor++
				}
			}
			return t, nil
		}

		// If not editing, only handle 'e' to start editing
		if msg.String() == "e" {
			t.editing = true
			configJSON, err := json.MarshalIndent(t.topic.Config, "", "  ")
			if err != nil {
				t.error = err
				return t, nil
			}
			t.jsonInput = string(configJSON)
			t.cursor = len(t.jsonInput)
		}
	}
	return t, nil
}

func (t *TopicDetails) View() string {
	if t.error != nil {
		return fmt.Sprintf("Error: %s", t.error)
	}

	var sb strings.Builder

	// Topic overview
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render("Topic Details\n\n"))
	sb.WriteString(fmt.Sprintf("Name: %s\n", t.topic.Name))
	sb.WriteString(fmt.Sprintf("Partitions: %d\n", t.topic.Partitions))
	sb.WriteString(fmt.Sprintf("Replication Factor: %d\n", t.topic.Replication))
	sb.WriteString(fmt.Sprintf("Internal: %v\n\n", t.topic.IsInternal))

	// Configuration
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render("Configuration\n\n"))

	if t.editing {
		sb.WriteString("Editing mode (press Enter to save, Esc to cancel):\n\n")
		sb.WriteString(t.jsonInput[:t.cursor])
		sb.WriteString("█") // Cursor
		sb.WriteString(t.jsonInput[t.cursor:])
	} else {
		if len(t.topic.Config) == 0 {
			sb.WriteString("No configuration available\n")
		} else {
			configJSON, err := json.MarshalIndent(t.topic.Config, "", "  ")
			if err != nil {
				sb.WriteString(fmt.Sprintf("Error formatting config: %s\n", err))
			} else {
				sb.WriteString(string(configJSON))
				sb.WriteString("\n")
			}
		}
		sb.WriteString("\nPress 'e' to edit configuration")
	}

	return sb.String()
}

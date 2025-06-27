package context

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/config"
	"github.com/snakeice/kafkalypse/internal/kafka"
	"github.com/snakeice/kafkalypse/internal/tui/messages"
	"github.com/snakeice/kafkalypse/internal/tui/styles"
	"github.com/spf13/viper"
)

// Brokers  string `yaml:"brokers"`
// Username string `yaml:"username"`
// Password string `yaml:"password"`

// SchemaRegistryUrl string `yaml:"schemaRegistryUrl"`
// RegistryUsername  string `yaml:"registryUsername"`
// RegistryPassword  string `yaml:"registryPassword"`

type NavigateToListMsg struct{}

type CreateContext struct {
	config *config.Configuration
	error  error

	fields     map[string]string
	fieldOrder []string
	cursor     int
	editing    bool
	editBuffer string

	sz messages.SizeMsg
}

func NewCreateContext(cfg *config.Configuration) *CreateContext {
	c := &CreateContext{
		config:     cfg,
		fields:     make(map[string]string),
		fieldOrder: []string{"name", "bootstrapServers", "securityProtocol"},
	}

	// Set default values
	c.fields["name"] = ""
	c.fields["bootstrapServers"] = "localhost:9093"
	c.fields["securityProtocol"] = "PLAINTEXT"

	return c
}

func (c *CreateContext) Init() tea.Cmd {
	return nil
}

func (c *CreateContext) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.SizeMsg:
		c.sz = msg

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+s":

			if c.fields["name"] == "" {
				c.error = fmt.Errorf("context name cannot be empty")
				return c, nil
			}

			// Save context
			ctx := &config.KafkaContext{
				BootstrapServers: c.fields["bootstrapServers"],
				SecurityProtocol: c.fields["securityProtocol"],
			}

			c.config.Contexts[c.fields["name"]] = ctx
			c.config.CurrentContext = c.fields["name"]

			viper.Set("contexts", c.config.Contexts)
			viper.Set("currentContext", c.config.CurrentContext)
			if err := c.config.SaveConfig(); err != nil {
				c.error = err
				return c, nil
			}

			return c, tea.Batch(
				messages.NavigateTo("contexts", true),
				kafka.Connect(ctx),
			)

		case "esc":
			if c.editing {
				c.editing = false
				c.editBuffer = ""
			} else {
				return c, messages.NavigateTo("contexts", true)
			}

		case "up":
			if !c.editing && c.cursor > 0 {
				c.cursor--
			}

		case "down":
			if !c.editing && c.cursor < len(c.fieldOrder)-1 {
				c.cursor++
			}

		case "enter":
			if !c.editing {
				c.editing = true
				c.editBuffer = c.fields[c.fieldOrder[c.cursor]]
			} else {
				c.fields[c.fieldOrder[c.cursor]] = c.editBuffer
				c.editing = false
				c.editBuffer = ""
			}

		case "backspace":
			if c.editing && len(c.editBuffer) > 0 {
				c.editBuffer = c.editBuffer[:len(c.editBuffer)-1]
			}

		default:
			if c.editing && len(msg.String()) == 1 {
				c.editBuffer += msg.String()
			}
		}
	}

	return c, nil
}

func (c *CreateContext) View() string {
	var content string

	// Error message
	if c.error != nil {
		content += styles.BasicStyle.
			Foreground(lipgloss.Color("#FF0000")).
			Render(fmt.Sprintf("Error: %s\n\n", c.error))
	}

	// Create title with border similar to content component
	title := "Create New Context"
	border := "╭"

	center := (c.sz.Width - 2 - lipgloss.Width(title)) / 2

	for range center {
		border += "─"
	}

	border += " " + title + " "

	for i := lipgloss.Width(border); i < c.sz.Width+1; i++ {
		border += "─"
	}

	border += "╮"
	content += styles.BasicStyle.Render(border) + "\n\n"

	// Form fields
	formStyle := styles.BasicStyle.
		Width(c.sz.Width - 8).
		MarginLeft(4)

	for i, field := range c.fieldOrder {
		// Highlight current field
		style := styles.BasicStyle
		if i == c.cursor {
			style = styles.TableHeaderSelected
		}

		// Show field name
		fieldContent := style.Render(fmt.Sprintf("%s: ", field))

		// Show field value or edit buffer
		valueStyle := styles.BasicStyle
		if c.editing && i == c.cursor {
			fieldContent += valueStyle.Render(c.editBuffer + "█")
		} else {
			fieldContent += valueStyle.Render(c.fields[field])
		}

		content += formStyle.Render(fieldContent) + "\n"
	}

	// Help text
	helpStyle := styles.BasicStyle.
		Faint(true).
		Width(c.sz.Width - 4).
		Align(lipgloss.Center)

	content += "\n"
	content += helpStyle.Render("Press Enter to edit field") + "\n"
	content += helpStyle.Render("Press Ctrl+S to save context") + "\n"
	content += helpStyle.Render("Press Esc to cancel")

	// Create bottom border
	bottomBorder := "╰"
	for i := 1; i < c.sz.Width; i++ {
		bottomBorder += "─"
	}
	bottomBorder += "╯"
	content += "\n\n" + styles.BasicStyle.Render(bottomBorder)

	// Create full-screen container using TableStyle similar to other components
	return styles.TableStyle.
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		Width(c.sz.Width).
		Height(c.sz.Height).
		Render(content)
}

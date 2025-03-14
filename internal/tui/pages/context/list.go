package context

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/config"
	"github.com/snakeice/kafkalypse/internal/constants"
	"github.com/snakeice/kafkalypse/internal/tui/components/shortcuts"
	"github.com/snakeice/kafkalypse/internal/tui/messages"
	"github.com/snakeice/kafkalypse/internal/tui/styles"
)

type ContextList struct {
	config    *config.Configuration
	error     error
	shortcuts *shortcuts.Shortcuts

	contexts []string
	cursor   int
	offset   int // For scrolling
}

func NewContextList(config *config.Configuration) *ContextList {
	c := &ContextList{
		config:    config,
		shortcuts: shortcuts.New(),
	}

	// Add shortcuts
	c.shortcuts.Add(shortcuts.Action{
		Key:         "n",
		Description: "Press 'n' to create new context",
		Handler: func() tea.Cmd {
			return messages.NavigateTo("create-context", true)
		},
	})

	c.shortcuts.Add(shortcuts.Action{
		Key:         "enter",
		Description: "Press Enter to select context",
		Handler: func() tea.Cmd {
			if len(c.contexts) > 0 {
				c.config.CurrentContext = c.contexts[c.cursor]
				if err := c.config.SaveConfig(); err != nil {
					c.error = err
					return nil
				}
				return messages.NavigateTo("topics", true)
			}
			return nil
		},
	})

	c.shortcuts.Add(shortcuts.Action{
		Key:         "pageup",
		Description: "PageUp to scroll up",
		Handler: func() tea.Cmd {
			c.cursor -= c.maxVisibleItems() - 2
			if c.cursor < 0 {
				c.cursor = 0
			}
			c.offset -= c.maxVisibleItems() - 2
			if c.offset < 0 {
				c.offset = 0
			}
			return nil
		},
	})

	c.shortcuts.Add(shortcuts.Action{
		Key:         "pagedown",
		Description: "PageDown to scroll down",
		Handler: func() tea.Cmd {
			c.cursor += c.maxVisibleItems() - 2
			if c.cursor >= len(c.contexts) {
				c.cursor = len(c.contexts) - 1
			}
			c.offset += c.maxVisibleItems() - 2
			if c.offset > len(c.contexts)-(c.maxVisibleItems()-2) {
				c.offset = len(c.contexts) - (c.maxVisibleItems() - 2)
			}
			if c.offset < 0 {
				c.offset = 0
			}
			return nil
		},
	})

	c.shortcuts.Add(shortcuts.Action{
		Key:         "home",
		Description: "Home to jump to start",
		Handler: func() tea.Cmd {
			c.cursor = 0
			c.offset = 0
			return nil
		},
	})

	c.shortcuts.Add(shortcuts.Action{
		Key:         "end",
		Description: "End to jump to end",
		Handler: func() tea.Cmd {
			c.cursor = len(c.contexts) - 1
			if len(c.contexts) > c.maxVisibleItems()-2 {
				c.offset = len(c.contexts) - (c.maxVisibleItems() - 2)
			}
			return nil
		},
	})

	return c
}

func (c *ContextList) Init() tea.Cmd {
	c.updateContextList()
	return nil
}

func (c *ContextList) updateContextList() {
	c.contexts = make([]string, 0, len(c.config.Contexts))
	for name := range c.config.Contexts {
		c.contexts = append(c.contexts, name)
	}
	sort.Strings(c.contexts)
}

func (c *ContextList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		// Handle navigation keys
		switch msg.String() {
		case "up":
			if c.cursor > 0 {
				c.cursor--
				// Scroll up if cursor is above viewport
				if c.cursor < c.offset {
					c.offset = c.cursor
				}
			}
			return c, nil
		case "down":
			if c.cursor < len(c.contexts)-1 {
				c.cursor++
				// Scroll down if cursor is below viewport
				if c.cursor >= c.offset+c.maxVisibleItems()-2 { // -2 for header and help text
					c.offset = c.cursor - c.maxVisibleItems() + 3
				}
			}
			return c, nil
		}

		// Handle other shortcuts
		if cmd := c.shortcuts.Handle(msg); cmd != nil {
			return c, cmd
		}
	}
	return c, nil
}

func (c *ContextList) maxVisibleItems() int {
	return constants.WindowHeight - 4 // Account for header, help text, and margins
}

func (c *ContextList) View() string {
	var content string

	// Error message
	if c.error != nil {
		content += styles.BasicStyle.
			Foreground(lipgloss.Color("#FF0000")).
			Render(fmt.Sprintf("Error: %s\n\n", c.error))
	}

	// Create title with border similar to content component
	title := "Kafka Contexts"
	border := "╭"

	center := (constants.WindowWidth - 2 - lipgloss.Width(title)) / 2

	for range center {
		border += "─"
	}

	border += " " + title + " "

	for i := lipgloss.Width(border); i < constants.WindowWidth+1; i++ {
		border += "─"
	}

	border += "╮"
	content += styles.BasicStyle.Render(border) + "\n\n"

	if len(c.contexts) == 0 {
		content += styles.BasicStyle.
			Width(constants.WindowWidth - 4).
			Align(lipgloss.Center).
			Render("No contexts found. Press 'n' to create a new context.")
	} else {
		// Calculate visible range
		endIdx := min(c.offset+c.maxVisibleItems()-2, len(c.contexts))

		// Context list
		listStyle := styles.BasicStyle.
			Width(constants.WindowWidth - 8).
			MarginLeft(4)

		for i := c.offset; i < endIdx; i++ {
			name := c.contexts[i]
			style := styles.BasicStyle
			prefix := "  "

			if i == c.cursor {
				style = styles.TableHeaderSelected
				prefix = "> "
			}
			if name == c.config.CurrentContext {
				style = style.Bold(true).Foreground(lipgloss.Color("#FFFF00"))
				prefix = "* "
			}
			if i == c.cursor && name == c.config.CurrentContext {
				prefix = ">*"
			}

			ctx := c.config.Contexts[name]
			content += listStyle.Render(style.Render(fmt.Sprintf("%s%s (%s)", prefix, name, ctx.BootstrapServers))) + "\n"
		}

		// Scroll indicators
		if c.offset > 0 {
			content = styles.BasicStyle.Align(lipgloss.Center).Render("↑") + "\n" + content
		}
		if endIdx < len(c.contexts) {
			content += styles.BasicStyle.Align(lipgloss.Center).Render("↓") + "\n"
		}
	}

	// Help text
	helpStyle := styles.BasicStyle.
		Faint(true).
		Width(constants.WindowWidth - 4).
		Align(lipgloss.Center)
	content += "\n" + helpStyle.Render(c.shortcuts.GetHelpText())

	// Create bottom border
	bottomBorder := "╰"
	for i := 1; i < constants.WindowWidth; i++ {
		bottomBorder += "─"
	}
	bottomBorder += "╯"
	content += "\n" + styles.BasicStyle.Render(bottomBorder)

	// Create full-screen container using TableStyle similar to other components
	return styles.TableStyle.
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		Width(constants.WindowWidth).
		Height(constants.WindowHeight).
		Render(content)
}

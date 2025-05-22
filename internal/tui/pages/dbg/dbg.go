package dbg

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/tui/messages"
)

type Dbg struct {
	sz messages.SizeMsg
}

func NewDbg() *Dbg {
	return &Dbg{}
}

func (c Dbg) Init() tea.Cmd { return nil }
func (c Dbg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.SizeMsg:
		c.sz = msg
	}

	return c, nil
}
func (c Dbg) View() string {
	v := "^\n"
	v += strings.Repeat("|\n", c.sz.Height-2)
	v += fmt.Sprintf("\nv ---> %d", c.sz.Height)
	return v

}

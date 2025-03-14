package tools

import tea "github.com/charmbracelet/bubbletea"

func WrapCmd(cmd tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return cmd
	}
}

package pagesstack

import tea "github.com/charmbracelet/bubbletea"

type PagesStack struct {
	pages []Page
}

func NewPagesStack() *PagesStack {
	return &PagesStack{}
}

func (s *PagesStack) Push(p Page) {
	s.pages = append(s.pages, p)
	s.pages[len(s.pages)-1].Active()

	if len(s.pages) > 1 {
		s.pages[len(s.pages)-2].Deactivate()
	}
}

func (s *PagesStack) Pop() {
	if len(s.pages) == 0 {
		return
	}

	s.pages = s.pages[:len(s.pages)-1]

	if len(s.pages) > 0 {
		s.pages[len(s.pages)-1].Active()
	}
}

func (s *PagesStack) Update(msg tea.Msg) tea.Cmd {
	if len(s.pages) == 0 {
		return nil
	}

	p := s.pages[len(s.pages)-1]
	newP, cmd := p.Update(msg)
	//nolint:sloppyTypeAssert
	if page, ok := newP.(Page); ok {
		s.pages[len(s.pages)-1] = page
	}

	return cmd
}

func (s *PagesStack) View() string {
	if len(s.pages) == 0 {
		return ""
	}

	return s.pages[len(s.pages)-1].View()
}

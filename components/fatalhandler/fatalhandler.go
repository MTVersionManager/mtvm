package fatalhandler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Model struct {
	err error
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if msg, ok := msg.(error); ok {
		m.err = msg
		return m, tea.Quit
	}
	return m, nil
}

func Handle(m Model) {
	if m.err != nil {
		log.Fatal(m.err)
	}
}

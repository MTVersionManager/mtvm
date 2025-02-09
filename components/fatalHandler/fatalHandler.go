package fatalHandler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Model struct {
	err error
}

// Update should only be called if there is an error.
func (m Model) Update(err error) (Model, tea.Cmd) {
	m.err = err
	return m, tea.Quit
}

func Handle(m Model) {
	if m.err != nil {
		log.Fatal(m.err)
	}
}

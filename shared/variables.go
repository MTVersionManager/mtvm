package shared

import (
	"github.com/MTVersionManager/mtvm/config"
	"github.com/charmbracelet/lipgloss"
)

var (
	Configuration config.Config
	CheckMark     = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).SetString("✓").String()
)

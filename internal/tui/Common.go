package tui

import "github.com/charmbracelet/lipgloss"

var (
	checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

type NewCountMsg struct {
	newCount int
}

type NewPacketReceivedMsg struct {
	value string
	count int
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

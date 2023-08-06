package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

type SpinnerModel struct {
	spinner        spinner.Model
	count          int
	startTimestamp string
	packetPreview  <-chan string
}

func (sm SpinnerModel) GetCount() int {
	return sm.count
}

func (sm SpinnerModel) Init() tea.Cmd {
	return tea.Batch(
		sm.spinner.Tick,
		sm.getPackets(),
	)
}

func (sm SpinnerModel) View() string {
	return lipgloss.JoinHorizontal(
		0.0,
		sm.spinner.View(),
		fmt.Sprintf(
			" Currently at %d packets retrived",
			sm.count,
		),
	)
}

func (sm SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return sm, tea.Quit
		}
	case NewPacketReceivedMsg:
		sm.count++
		return sm, tea.Batch(sm.getPackets(), tea.Printf("%s %s", checkMark, msg.value))
	case spinner.TickMsg:
		var cmd tea.Cmd
		sm.spinner, cmd = sm.spinner.Update(msg)
		return sm, cmd
	}

	return sm, nil
}

func NewPacketInfinitSpinner(pckPreview <-chan string) *SpinnerModel {
	return &SpinnerModel{
		spinner: spinner.New(
			spinner.WithStyle(
				lipgloss.NewStyle().
					Foreground(lipgloss.Color("#ff6788")),
			),
		),
		count:          0,
		packetPreview:  pckPreview,
		startTimestamp: time.Now().String(),
	}
}

func (sm SpinnerModel) getPackets() tea.Cmd {
	return func() tea.Msg {
		return NewPacketReceivedMsg{
			value: <-sm.packetPreview,
		}
	}
}

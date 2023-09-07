package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"gosniff/internal/database"
)

type PacketViewport struct {
	Viewport viewport.Model
}

func (v PacketViewport) View() string {
	return v.Viewport.View()
}

func (v PacketViewport) Init() tea.Cmd {
	return tea.Batch(v.Viewport.Init())
}

func (v PacketViewport) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch _msg := msg.(type) {
	case tea.KeyMsg:
		switch _msg.String() {
		case "q", "shift+q", "ctrl+c":
			return v, tea.Quit
		}
	case tea.WindowSizeMsg:
		v.Viewport = viewport.New(_msg.Width, _msg.Height)
		packets, _ := database.FetchLastInsertedPackets()
		v.Viewport.SetContent(packets[0].Content)
		return v, nil
	}
	var cmd tea.Cmd
	v.Viewport, cmd = v.Viewport.Update(msg)

	return v, cmd
}

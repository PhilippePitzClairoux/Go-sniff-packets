package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gosniff/internal/database"
	"gosniff/internal/internal"
)

func NewBrowsePacketTui() *PacketsMainMenu {
	items := make([]list.Item, 0)
	lst := list.New(items, list.NewDefaultDelegate(), 0, 0)

	lst.Title = "SUPER DUPER PACKET SEARCHER"

	return &PacketsMainMenu{
		viewport: PacketViewport{},
		list:     lst,
	}
}

type PacketsMainMenu struct {
	// add various menu interactions
	viewport     PacketViewport // view single packet
	list         list.Model     // search for packets and display them (via viewport)
	showViewport bool
	items        []list.Item
}

type newPacketReceivedMsg struct {
	packets []list.Item
}

func (m PacketsMainMenu) Init() tea.Cmd {
	return func() tea.Msg {
		packets, err := database.FetchLastInsertedPackets()
		if err != nil {
			return nil
		}

		var items []list.Item
		for _, packet := range packets {
			items = append(items, packet)
		}

		return newPacketReceivedMsg{
			packets: items,
		}
	}
}

func (m PacketsMainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q", "ctrl+c":
			if m.showViewport {
				m.showViewport = false
				return m, nil
			}
			return m, tea.Quit
		case "f2":
			//create a new search box to load data from database
		case "enter":
			if !m.showViewport && m.items != nil {
				m.showViewport = true
				return m, nil
			}

		}
	case newPacketReceivedMsg:
		m.items = append(m.items, msg.packets...)
		m.list.SetItems(m.items)
		return m, m.Init()
	case tea.WindowSizeMsg:
		m.list.SetHeight(msg.Height)
		m.list.SetWidth(msg.Width)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m PacketsMainMenu) View() string {
	if m.showViewport && m.items != nil {
		m.viewport.Viewport.SetContent(m.list.SelectedItem().(internal.Packet).Content)
		return m.viewport.View()
	}

	return m.list.View()
}

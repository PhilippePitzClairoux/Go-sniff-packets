package tui

// TODO rework thise whole code (removed an event so might have to either reuse PacketInfinitSpinner event or just create a new one)
//
//import (
//	"fmt"
//	"github.com/charmbracelet/bubbles/progress"
//	"github.com/charmbracelet/bubbles/spinner"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/lipgloss"
//	"strconv"
//	"strings"
//	"time"
//)
//
//type InitialLoader struct {
//	PacketCountAware
//	width             int
//	height            int
//	currentAmount     int
//	lastCurrentAmount int
//	spinner           spinner.Model
//	progress          progress.Model
//	done              bool
//}
//
//var (
//	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
//	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
//
//)
//
//const threshHold = 250
//
//func NewInitialLoader() *InitialLoader {
//	p := progress.New(
//		progress.WithDefaultGradient(),
//		progress.WithWidth(40),
//		progress.WithoutPercentage(),
//	)
//	s := spinner.New()
//	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
//	return &InitialLoader{
//		spinner:  s,
//		progress: p,
//		PacketCountAware: PacketCountAware{
//			startTimestamp: time.Now().String(),
//		},
//	}
//}
//
//func (m InitialLoader) Init() tea.Cmd {
//	return tea.Batch(m.updatePacketCount())
//}
//
//func (m InitialLoader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.WindowSizeMsg:
//		m.width, m.height = msg.Width, msg.Height
//	case tea.KeyMsg:
//		switch msg.String() {
//		case "ctrl+c", "esc", "q":
//			return m, tea.Quit
//		}
//	case NewCountMsg:
//		if msg.newCount >= threshHold {
//			// Everything's been installed. We're done!
//			m.done = true
//			return m, tea.Quit
//		}
//
//		if m.lastCurrentAmount < msg.newCount {
//			m.lastCurrentAmount = m.currentAmount
//			m.currentAmount = msg.newCount
//
//			// Update progress bar
//			progressCmd := m.progress.SetPercent(float64(m.currentAmount) / float64(threshHold))
//
//			return m, tea.Batch(
//				progressCmd,
//				m.updatePacketCount(), // update internal count
//			)
//		}
//
//		return m, m.updatePacketCount()
//	case spinner.TickMsg:
//		var cmd tea.Cmd
//		m.spinner, cmd = m.spinner.Update(msg)
//		return m, cmd
//	case progress.FrameMsg:
//		newModel, cmd := m.progress.Update(msg)
//		if newModel, ok := newModel.(progress.Model); ok {
//			m.progress = newModel
//		}
//		return m, cmd
//	}
//	return m, nil
//}
//
//func (m InitialLoader) View() string {
//	w := lipgloss.Width(fmt.Sprintf("%d", threshHold))
//
//	if m.done {
//		return doneStyle.Render(fmt.Sprintf("Done! Finished getting the initial %d packets!\n", threshHold))
//	}
//
//	pckCount := fmt.Sprintf(" %*d/%*d", w, m.currentAmount, w, threshHold)
//
//	spin := m.spinner.View() + " "
//	prog := m.progress.View()
//	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pckCount))
//
//	currPckCount := currentPkgNameStyle.Render("Current amount of captured packets : " + strconv.Itoa(m.currentAmount))
//	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render(currPckCount)
//
//	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pckCount))
//	gap := strings.Repeat(" ", cellsRemaining)
//
//	return spin + info + gap + prog + pckCount
//}

package tui

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"Qischer/player-tui/internal/player"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
)

type model struct {
	status int
	err    error
	width  int
	height int
	state  player.PlayerState

	// components
	progress progress.Model
	// Styles
	styles *Styles

	// server
	last int64
	quit chan struct{}

	// Prompt
	isPrompt bool
	prompt   Prompt
}

// Model messages
type statusMsg struct {
	code  int
	state player.PlayerState
}

type Prompt struct {
	message string
	url     string
}

type (
	waitMsg struct{}
	errMsg  struct{ err error }
)

// Styling
type Styles struct {
	BorderColor lipgloss.Color
	PlayerBox   lipgloss.Style
}

const (
	padding = 2
)

func (e errMsg) Error() string { return e.err.Error() }

func DefaultStyles() *Styles {
	s := &Styles{}
	s.BorderColor = lipgloss.Color("#FF7CCB")
	s.PlayerBox = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		Padding(2).Width(80)
	return s
}

func (m model) Init() tea.Cmd {
	go StartServer(m.quit)

	if os.Getenv("SPOTIFY_REFRESH_TOKEN") == "" {
		return getAuthLink
	}
	return updatePlayerState(0)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case statusMsg:
		m.status = msg.code
		m.state = msg.state
		m.last = time.Now().UnixMilli()
		return m, updatePlayerState(m.last)

	case waitMsg:
		return m, updatePlayerState(m.last)

	case errMsg:
		m.err = msg

	case Prompt:
		log.Println("prompt received")
		m.isPrompt = true
		m.prompt = msg

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			close(m.quit)
			return m, tea.Quit
		case " ":
			// log.Println("Play/Pause")
			return m, togglePlayback(m.state.IsPlaying, m.last)

		case "n":
			log.Println("Next")
			return m, updatePlayerState(m.last)

		case "p":
			log.Println("Prev")
			return m, updatePlayerState(m.last)

		case "y":
			log.Println("Accepted Prompt")
			m.AcceptPrompt()
			m.isPrompt = false
			return m, updatePlayerState(m.last)
		}
	}

	return m, nil
}

func (m model) AcceptPrompt() {
	if err := browser.OpenURL(m.prompt.url); err != nil {
		log.Fatal(err)
	}

	for os.Getenv("SPOTIFY_REFRESH_TOKEN") == "" {
	}
}

func (m model) View() string {
	m.styles = DefaultStyles()
	pad := strings.Repeat(" ", padding)

	if m.err != nil {
		return "An error occurred"
	}

	ui := ""

	if m.isPrompt {
		m := lipgloss.JoinHorizontal(lipgloss.Center, m.prompt.message)
		btn := lipgloss.JoinHorizontal(lipgloss.Center, "Redirect (Y)", "Quit (N)")

		ui = lipgloss.JoinVertical(lipgloss.Center, m, btn)
	} else if m.status == http.StatusNoContent {
		ui = "Player not active"
	} else {
		a := printArtists(m.state.Item.Artists)
		s := fmt.Sprintf("\033[1m%v\033[0m\n\033[90m%v\033[0m\n", m.state.Item.Name, a)

		end := parseTime(int(m.state.Item.Duration))
		cur := parseTime(int(m.state.ProgressMS))
		bar := m.progress.ViewAs(float64(m.state.ProgressMS) / float64(m.state.Item.Duration))

		// time
		ts := lipgloss.JoinHorizontal(lipgloss.Center, cur, pad, bar, pad, end)
		ui = lipgloss.JoinVertical(lipgloss.Center, s, ts)
	}
	return lipgloss.Place(m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.styles.PlayerBox.Render(ui),
	)
}

func NewModel(quit chan struct{}) model {
	styles := DefaultStyles()

	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	prog.ShowPercentage = false
	prog.Width = 60

	return model{styles: styles, progress: prog, quit: quit}
}

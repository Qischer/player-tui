package tui

import (
	"Qischer/player-tui/internal/player"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
  status  int 
  err     error
  width   int
  height  int
  state   player.PlayerState

  //components
  progress progress.Model
  //Styles 
  styles  *Styles

  //server 
  last    int64
  quit    chan struct{}
}

//Model messages
type statusMsg struct{
  code  int
  state player.PlayerState
}

type waitMsg struct{}
type errMsg struct{err error}

//Styling
type Styles struct {
  BorderColor lipgloss.Color
  PlayerBox   lipgloss.Style
}

const (
  padding = 2;
)

func (e errMsg) Error() string { return e.err.Error() }

func DefaultStyles() *Styles {
  s := &Styles{}
  s.BorderColor = lipgloss.Color("201")
  s.PlayerBox = lipgloss.NewStyle().
          BorderForeground(s.BorderColor).
          BorderStyle(lipgloss.RoundedBorder()).
          Align(lipgloss.Center).
          Padding(2).Width(80)
  return s
}

func (m model) Init() (tea.Cmd) {
  go startServer(m.quit)
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
    return m, nil

  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c","q":
      close(m.quit)
      return m, tea.Quit
    }
  }

  return m, nil
}

func (m model) View() string {
  m.styles = DefaultStyles()
  pad := strings.Repeat(" ", padding)

  if m.err != nil { 
    return "An error occurred"
  }

  if m.status == http.StatusNoContent {
    return "Player not active"
  }

  a := printArtists(m.state.Item.Artists)
  s := fmt.Sprintf("\033[1m%v\033[0m\n\033[90m%v\033[0m\n", m.state.Item.Name, a)

  end := parseTime(int(m.state.Item.Duration))
  cur := parseTime(int(m.state.ProgressMS)) 
  bar := m.progress.ViewAs(float64(m.state.ProgressMS) / float64(m.state.Item.Duration))

  //time
  ts := lipgloss.JoinHorizontal(lipgloss.Center, cur, pad, bar, pad, end)
  vui := lipgloss.JoinVertical(lipgloss.Center, s, ts)
  hui := lipgloss.JoinHorizontal(lipgloss.Center, vui)
  return lipgloss.Place(m.width, 
      m.height, 
      lipgloss.Center, 
      lipgloss.Center, 
      m.styles.PlayerBox.Render(hui),
  )
}

func NewModel(quit chan struct{}) model{ 
  styles := DefaultStyles()
  
  prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
  prog.ShowPercentage = false
  prog.Width = 60

  return model{styles: styles, progress: prog, quit: quit}
}

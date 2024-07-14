package tui

import (
	"Qischer/player-tui/internal/player"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const url = "http://localhost:6969/player/state"

type model struct {
  status  int 
  err     error
  width   int
  height  int
  state   player.PlayerState

  //Styles 
  styles  *Styles
}

//Model messages
type statusMsg struct{
  code  int
  state player.PlayerState
}

type errMsg struct{err error}

//Styling
type Styles struct {
  BorderColor lipgloss.Color
  PlayerBox   lipgloss.Style
}

func (e errMsg) Error() string { return e.err.Error() }

func DefaultStyles() *Styles {
  s := &Styles{}
  s.BorderColor = lipgloss.Color("201")
  s.PlayerBox = lipgloss.NewStyle().
        BorderForeground(s.BorderColor).
        BorderStyle(lipgloss.RoundedBorder()).
        Padding(3).Width(80)
  return s
}

func (m model) Init() (tea.Cmd) {
  return reqState
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height

  case statusMsg:
    m.status = msg.code
    m.state = msg.state
    return m, reqState

  case errMsg: m.err = msg
    return m, nil

  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c","q":
      return m, tea.Quit
    }
  }
  return m, reqState
}

func (m model) View() string {
  m.styles = DefaultStyles()

  if m.err != nil { 
    return "An error occurred"
  }

  if m.status == http.StatusNoContent {
    return "Player not active"
  }

  s := fmt.Sprintf("Playing : %v\n", m.state.Item.Name)

  end := parseTime(int(m.state.Item.Duration))
  cur := parseTime(int(m.state.ProgressMS)) 

  //time
  ts := lipgloss.JoinHorizontal(lipgloss.Center, cur, "[status bar]", end)
  ui := lipgloss.JoinVertical(lipgloss.Center, s, ts)
  return lipgloss.Place(m.width, 
          m.height, 
          lipgloss.Center, 
          lipgloss.Center, 
          m.styles.PlayerBox.Render(ui),
        )
  }

func NewModel() model{
  styles := DefaultStyles()
  return model{styles: styles}
}

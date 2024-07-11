package tui

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "http://localhost:6969/player/state"

type model struct {
  status  int 
  err     error
}

type statusMsg int

type errMsg struct{err error}

func (e errMsg) Error() string { return e.err.Error() }

func checkServer() tea.Msg {

  c := &http.Client{Timeout: 10 * time.Second}
  res, err := c.Get(url)

  if err != nil {
    return errMsg{err}
  }

  return statusMsg(res.StatusCode)
}

func (m model) Init() (tea.Cmd) {
  
  return checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
    
  case statusMsg:
    m.status = int(msg)
    return m, checkServer

  case errMsg: 
    m.err = msg
    return m, nil

  }
  return m, checkServer
}

func (m model) View() string {
  if m.err != nil {
    return fmt.Sprintf("\nError encountered: %v\n\n", m.err)
  }

  s := fmt.Sprintf("Checking %s...", url)

  if m.status > 0 { 
    s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
  }

  return "\n" + s + "\n\n"
}

func main() {
  if _, err := tea.NewProgram(model{}, tea.WithAltScreen()).Run(); err != nil {
    fmt.Println("Something went wrong")
    os.Exit(1)
  }
}

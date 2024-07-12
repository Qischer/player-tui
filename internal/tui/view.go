package tui

import (
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

func (m model) Init() (tea.Cmd) {
  
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
    
  case statusMsg:
    m.status = int(msg)
    return m, nil

  case errMsg: 
    m.err = msg
    return m, nil

  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c","q":
      return m, tea.Quit
    }
  }
  return m, nil
}

func (m model) View() string {
  return "Hello World!"
}

func NewModel() model{
  return model{}
}

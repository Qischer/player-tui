package tui

import (
	"Qischer/player-tui/internal/player"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "http://localhost:6969/player/state"

type model struct {
  status  int 
  err     error
  state   player.PlayerState
}

type statusMsg struct{
  code  int
  state player.PlayerState
}

type errMsg struct{err error}

func (e errMsg) Error() string { return e.err.Error() }

func reqState() tea.Msg {
  //http client
  c := &http.Client{}
  res, err := c.Get(url)
  if err != nil {
    log.Println(err)
    return errMsg{err}
  }

  //get res body 
  state := &player.PlayerState{}
  if res.StatusCode == http.StatusOK {
    if e:= json.NewDecoder(res.Body).Decode(state); e != nil {
      log.Fatal(e)
    }
  }
  
  return statusMsg{
    code: res.StatusCode,
    state: *state,
  } 
}

func (m model) Init() (tea.Cmd) {
  
  return reqState
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
    
  case statusMsg:
    m.status = msg.code
    m.state = msg.state
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
  if m.err != nil { 
    return "An error occurred"
  }

  if m.status == http.StatusNoContent {
    return "Player not active"
  }

  s := fmt.Sprintf("Playing : %v\n", m.state.Item.Name)
  return s 
}

func NewModel() model{
  return model{}
}

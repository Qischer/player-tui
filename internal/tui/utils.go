package tui

import (
	"Qischer/player-tui/internal/player"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func parseTime(ms int) string {
  m := int((ms/ 1000) / 60)
  s := int((ms/ 1000) % 60)

  return fmt.Sprintf("%d:%2d", m, s)
}

func reqState() tea.Msg {
  //http Client
  time.Sleep(1 * time.Second)

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

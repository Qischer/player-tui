package tui

import (
	"Qischer/player-tui/internal/player"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
  "os"

	tea "github.com/charmbracelet/bubbletea"
	yaml "gopkg.in/yaml.v3"
)

const url = "http://localhost:6969/player/state"

func parseTime(ms int) string {
  m := int((ms/ 1000) / 60)
  s := int((ms/ 1000) % 60)

  return fmt.Sprintf("%d:%02d", m, s)
}

func startServer(q chan struct{}) {
  //Get access codes
  f, err := os.ReadFile(".env.yaml")
  if err != nil {
    log.Fatal(err)
  }

  m := player.SpotifyKeys{}
  err = yaml.Unmarshal(f, &m)
  if err != nil {
    log.Fatal(err)
  }

  os.Setenv("SPOTIFY_CLIENT_ID", m.ClientID)
  os.Setenv("SPOTIFY_CLIENT_SECRET", m.ClientSecret)
  os.Setenv("SPOTIFY_REFRESH_TOKEN", m.RefreshToken)

  router := http.NewServeMux()
  player.LoadRoutes(router)

  server := http.Server{
    Addr: ":6969",
    Handler: router,
  }

  log.Println("Listening on port 6969")
  if err := server.ListenAndServe(); err != nil {
    log.Println(err)
    close(q)
  }
}

func updatePlayerState(last int64) tea.Cmd {

  del := time.Now().UnixMilli() - last

  if del <= 1 {
    return func() tea.Msg {
      return waitMsg{}
    }
  }

  return func() tea.Msg {
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
}

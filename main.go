package main

import (
	"Qischer/player-tui/internal/player"
	"Qischer/player-tui/internal/tui"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	yaml "gopkg.in/yaml.v3"
)

func preclose() {
  m := &player.SpotifyKeys{
    ClientID: os.Getenv("SPOTIFY_CLIENT_ID"),
    ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
    RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
  }

  if m.ClientID == "" {
    log.Fatal("Premature Exit. Prevented api key clear")
  }

  out, err := yaml.Marshal(m)
  if err != nil {
    log.Fatal(err)
  }

  err = os.WriteFile(".env.yaml", out, 0644)
  if err != nil {
    log.Fatal(err)
  }

  log.Println("Saved to yaml file")
}

func main() { 
  quit := make(chan struct{})
  if len(os.Args) >= 2 && os.Args[1] == "-d" {
    tui.StartServer(quit)
    os.Exit(0)
  }

  f, err := os.OpenFile("log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  if err != nil {
    log.Fatal(err)
  }

  defer f.Close()
  log.SetOutput(f)

  p := tea.NewProgram(tui.NewModel(quit), tea.WithAltScreen())
  if _, err := p.Run(); err != nil {
    log.Println("Error running program:", err)
    close(quit) // Signal to stop the HTTP server
  }

  <-quit
  preclose()
  log.Println("Goodbye!")
  os.Exit(0)

}

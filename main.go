package main

import (
	"Qischer/player-tui/internal/player"
	"Qischer/player-tui/internal/tui"
	"fmt"
	"log"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	yaml "gopkg.in/yaml.v3"
)

func startServer(q chan struct{}) {
  fmt.Println("Run TUI")

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

func runTUI() {
  //tui.main()
}

func main() { 
  quit := make(chan struct{})
  f, err := os.OpenFile("log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  if err != nil {
    log.Fatal(err)
  }

  defer f.Close()

  log.SetOutput(f)

  go startServer(quit)
  // Create and start the Bubble Tea TUI
  p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())

  // Run the TUI in a separate goroutine
  go func() {
    if _, err := p.Run(); err != nil {
      log.Println("Error running program:", err)
      close(quit) // Signal to stop the HTTP server
    }
  }()
  <-quit

  preclose()
  log.Println("Goodbye!")
  os.Exit(0)

}

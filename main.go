package main

import (
  "Qischer/player-tui/internal/player"
  "net/http"
  "os"
  "os/signal"
  "syscall"

  log "github.com/charmbracelet/log"
  "gopkg.in/yaml.v3"
)

func run() {
  log.SetLevel(log.DebugLevel)
  log.Info("Run Player TUI")

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

  log.Info("Listening on port 6969")
  server.ListenAndServe()
}

func preclose() {
  m := &player.SpotifyKeys{
    ClientID: os.Getenv("SPOTIFY_CLIENT_ID"),
    ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
    RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
  }

  out, err := yaml.Marshal(m)
  if err != nil {
    log.Fatal(err)
  }

  err = os.WriteFile(".env.yaml", out, 0644)
  if err != nil {
    log.Fatal(err)
  }

  log.Debug("Saved to yaml file")
}

func main() { 

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)

  go run()

  <-c 

  preclose()
  log.Info("Goodbye!")
  os.Exit(0)

}

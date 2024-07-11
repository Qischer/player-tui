package main

import (
	"Qischer/player-tui/internal/player"
	"os"
	"os/signal"
	"syscall"

	log "github.com/charmbracelet/log"
	yaml "gopkg.in/yaml.v3"
)

func preClose() {
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

  log.Debug("Saved to yaml file")
}

func runTUI() {
  //tui.main()
}

func main() { 
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)

  go player.StartServer()
  //tui.Main()

  <-c 

  preClose()
  log.Info("Goodbye!")
  os.Exit(0)

}

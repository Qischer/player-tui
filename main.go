package main

import (
	"Qischer/player-tui/internal/tui"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func preclose() {
  // m := &player.SpotifyKeys{
  m := make(map[string]string)
  m["SPOTIFY_CLIENT_ID"] = os.Getenv("SPOTIFY_CLIENT_ID")
  m["SPOTIFY_CLIENT_SECRET"] = os.Getenv("SPOTIFY_CLIENT_SECRET")
  m["SPOTIFY_REFRESH_TOKEN"] = os.Getenv("SPOTIFY_REFRESH_TOKEN")

  err := godotenv.Write(m, ".env")
  if err != nil {
    log.Fatal(err)
  }
  log.Println("Updated env")
}

func main() { 
  if err := godotenv.Load(); err != nil {
    log.Fatal(err)
  }

  quit := make(chan struct{})
  if len(os.Args) >= 2 && os.Args[1] == "-d" {
    go tui.StartServer(quit)

  } else {
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
  }

  <-quit
  preclose()
  log.Println("Goodbye!")
  os.Exit(0)

}

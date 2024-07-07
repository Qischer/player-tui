package main

import (
	"Qischer/player-tui/internal/player"
	"net/http"

	log "github.com/charmbracelet/log"
)

func main() { 
  log.SetLevel(log.DebugLevel)
  log.Info("Run Player TUI")
  
  router := http.NewServeMux()
  player.LoadRoutes(router)

  server := http.Server{
    Addr: ":6969",
    Handler: router,
  }

  log.Info("Listening on port 6969")
  server.ListenAndServe()
 }

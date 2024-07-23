package main

import (
	"log"
	"os"
	"os/signal"

	"Qischer/player-tui/internal/tui"

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
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go tui.StartServer(quit)

		<-c
		close(quit)
	} else {
		f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
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

	preclose()
	log.Println("Goodbye!")
	os.Exit(0)
}

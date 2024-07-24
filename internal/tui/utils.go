package tui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"Qischer/player-tui/internal/player"

	tea "github.com/charmbracelet/bubbletea"
)

func parseTime(ms int) string {
	m := int((ms / 1000) / 60)
	s := int((ms / 1000) % 60)

	return fmt.Sprintf("%d:%02d", m, s)
}

func printArtists(a []player.Artist) string {
	sarr := []string{}
	for _, artist := range a {
		sarr = append(sarr, artist.Name)
	}

	return strings.Join(sarr, ", ")
}

func StartServer(q chan struct{}) {
	router := http.NewServeMux()
	player.LoadRoutes(router)

	server := http.Server{
		Addr:    ":6969",
		Handler: router,
	}

	log.Println("Listening on port 6969")
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
		close(q)
	}
}

func updatePlayerState(last int64) tea.Cmd {
	const u = "http://localhost:6969/player/state"
	del := time.Now().UnixMilli() - last

	if del <= 1 {
		return func() tea.Msg {
			return waitMsg{}
		}
	}

	return func() tea.Msg {
		// http Client
		time.Sleep(1 * time.Second)

		c := &http.Client{}
		res, err := c.Get(u)
		if err != nil {
			log.Println(err)
			return errMsg{err}
		}

		// get res body
		state := &player.PlayerState{}
		if res.StatusCode == http.StatusOK {
			if e := json.NewDecoder(res.Body).Decode(state); e != nil {
				log.Fatal(e)
			}
		}

		return statusMsg{
			code:  res.StatusCode,
			state: *state,
		}
	}
}

func (m model) togglePlayback() {
	var link string
	if m.state.IsPlaying {
		link = "http://localhost:6969/player/pause"
	} else {
		link = "http://localhost:6969/player/play"
	}

	up, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}

	req := &http.Request{
		Method: "PUT",
		URL:    up,
	}

	c := &http.Client{}
	_, er := c.Do(req)
	if er != nil {
		log.Fatal(er)
	}
}

func getAuthLink() tea.Msg {
	u := "http://localhost:6969/auth"
	res, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}

	auth := &player.PlayerAuth{}
	if err := json.NewDecoder(res.Body).Decode(auth); err != nil {
		log.Fatal(err)
	}

	return Prompt{
		message: "Redirect to Spotify for authentication?",
		url:     auth.URL,
	}
}

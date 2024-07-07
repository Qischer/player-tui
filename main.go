package main

import (
	"net/http"
	"net/url"

	log "github.com/charmbracelet/log"
)

var (
  CLIENT_ID = ""
  SCOPE = "user-read-playback-state"
  REDIRECT_URI = "http://localhost:6969/callback"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) {
  log.Info("Request:", r.Method, r.URL.Path)
  w.Write([]byte("You will be redirected to Spotify\n"))
  
  q := make(url.Values)
  q.Add("response_type", "code")
  q.Add("client_id", CLIENT_ID)
  q.Add("scope", SCOPE)
  q.Add("redirect_uri", REDIRECT_URI)
 
  u := &url.URL{
    Scheme: "https",
    Host: "accounts.spotify.com",
    Path: "authorize",
    RawQuery: q.Encode(),
  }
  
  w.Write([]byte(u.String()))
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
  log.Info("Request:", r.Method, r.URL.Path)
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Hi there, Listener!!"))
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
  code := r.URL.Query().Get("code")

  log.Info("code: ",code)
}

func main() {
  log.Info("Run Player TUI")

  router := http.NewServeMux()
  router.HandleFunc("GET /player", HandleIndex)
  router.HandleFunc("GET /auth", HandleAuth)
  router.HandleFunc("GET /callback", HandleCallback)

  server := http.Server{
    Addr: ":6969",
    Handler: router,
  }

  log.Info("Listening on port 6969")
  server.ListenAndServe()
 }

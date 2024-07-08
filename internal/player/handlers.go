package player

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	log "github.com/charmbracelet/log"
)

var (
  SCOPE = "user-read-playback-state"
  REDIRECT_URI = "http://localhost:6969/callback"
)

type Handlers struct {}

func (h *Handlers) HandleIndex(w http.ResponseWriter, r *http.Request) {
  log.Info("Request Made", "method", r.Method, "path", r.URL.Path)
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Hi there, Listener!!"))
}

func (h *Handlers) HandleAuth(w http.ResponseWriter, r *http.Request) {
  log.Info("Request:", r.Method, r.URL.Path)
  w.Write([]byte("You will be redirected to Spotify\n"))

  q := make(url.Values)
  q.Add("response_type", "code")
  q.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
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

func (h *Handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
  code := r.URL.Query().Get("code")

  req := &AccessRequest{
    GrantType: "authorization_code",
    AuthCode: code,
  }

  access, err := req.MakeRequest()
  if err != nil {
    log.Error(err)
  }
  
  os.Setenv("SPOTIFY_ACCESS_CODE", access.AccessToken)
  os.Setenv("SPOTIFY_REFRESH_TOKEN", access.RefreshToken)

  log.Debug("Check code", "real", access.AccessToken, "env", os.Getenv("SPOTIFY_ACCESS_CODE"))
  log.Info("Updated access_code and refresh_token to env")
}



func (h *Handlers) HandleGetState(w http.ResponseWriter, r *http.Request) {
  u := &url.URL{
    Scheme: "https",
    Host: "api.spotify.com",
    Path: "v1/me/player",
  }
  
  accessTok := os.Getenv("SPOTIFY_ACCESS_CODE")
  
  if accessTok == "" {
    req := &AccessRequest{
      GrantType: "refresh_token",
      RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
    } 

    access, err := req.MakeRequest()
    if err != nil { 
      log.Error(err)
    } 
    os.Setenv("SPOTIFY_ACCESS_CODE", access.AccessToken)

    log.Info("Updated access_code to env") 
  }

  log.Debug("Get player state", "access_code", os.Getenv("SPOTIFY_ACCESS_CODE")) 

  req, err := http.NewRequest("GET", u.String(), bytes.NewBuffer([]byte("")))
  if err != nil {
    log.Fatal(err)
  }

  req.Header.Add("Authorization", "Bearer " + os.Getenv("SPOTIFY_ACCESS_CODE"))

  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
    log.Fatal(err)
  }

  if res.StatusCode != http.StatusOK {
    errPost := ApiError{}
    e := json.NewDecoder(res.Body).Decode(&errPost)
    
    if e.Error() == "EOF" {
      w.WriteHeader(400)
      w.Write([]byte("No active player"))
      return
    }

    if e != nil {
      log.Fatal(e)
    }

    log.Error("Error Occured", "error", errPost)
    return
  }

  state := &PlayerState{}
  e := json.NewDecoder(res.Body).Decode(state)
  if e != nil {
    log.Error(e)
  }

  log.Info("Player State", "state", state)
  e = json.NewEncoder(w).Encode(state)
  if e != nil {
    log.Fatal(e)
  }
}




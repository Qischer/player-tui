package player

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
  SCOPE = "user-read-playback-state"
  REDIRECT_URI = "http://localhost:6969/callback"
)

type Handlers struct {}

func (h *Handlers) HandleIndex(w http.ResponseWriter, r *http.Request) {
  log.Println("Request Made", "method", r.Method, "path", r.URL.Path)
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Hi there, Listener!!"))
}

func (h *Handlers) HandleAuth(w http.ResponseWriter, r *http.Request) {
  log.Println("Request:", r.Method, r.URL.Path)
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
    log.Fatal(err)
  }
  
  os.Setenv("SPOTIFY_ACCESS_CODE", access.AccessToken)
  os.Setenv("SPOTIFY_REFRESH_TOKEN", access.RefreshToken)

  log.Println("Updated access_code and refresh_token to env")
}



func (h *Handlers) HandleGetState(w http.ResponseWriter, r *http.Request) {
  log.Println("Request player state")

  //last req time 
  // lt := r.URL.Query().Get("prev")

  u := &url.URL{
    Scheme: "https",
    Host: "api.spotify.com",
    Path: "v1/me/player",
  }
  
  atok := os.Getenv("SPOTIFY_ACCESS_CODE")
  
  if atok == "" {
    req := &AccessRequest{
      GrantType: "refresh_token",
      RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
    } 

    access, err := req.MakeRequest()
    if err != nil { 
      log.Fatal(err)
    } 
    os.Setenv("SPOTIFY_ACCESS_CODE", access.AccessToken)

    log.Println("Updated access_code to env") 
  }

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

  if res.StatusCode == http.StatusNoContent {
    w.WriteHeader(http.StatusNoContent)
    return
  }

  if res.StatusCode != http.StatusOK {
    errPost := ApiError{}
    e := json.NewDecoder(res.Body).Decode(&errPost)

    if e != nil {
      log.Fatal(e)
    }

    log.Fatal("Error Occured", "error", errPost)
    return
  }

  state := &PlayerState{}
  e := json.NewDecoder(res.Body).Decode(state)
  if e != nil {
    log.Fatal(e)
  }

  log.Println("Player State", "state", state)
  w.WriteHeader(http.StatusOK)
  e = json.NewEncoder(w).Encode(state)
  if e != nil {
    log.Fatal(e)
  }
}

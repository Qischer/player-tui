package player

import (
 	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
  "os"
 
  log "github.com/charmbracelet/log"
)

var (
  CLIENT_ID = os.Getenv("SPOTIFY_CLIENT_ID")
  CLIENT_SECRET = os.Getenv("SPOTIFY_CLIENT_SECRET")
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

func (h *Handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
  code := r.URL.Query().Get("code")
  
  u := &url.URL{
    Scheme: "https",
    Host: "accounts.spotify.com",
    Path: "api/token",
  }

  body := make(url.Values)
  body.Add("grant_type", "authorization_code")
  body.Add("code", code)
  body.Add("redirect_uri", REDIRECT_URI)

  r, err := http.NewRequest("POST", u.String(), bytes.NewBuffer([]byte(body.Encode())))
  if err != nil {
    log.Fatal(err)
  }
  
  id_secret := string(CLIENT_ID+":"+CLIENT_SECRET)
  auth_val := `Basic ` + base64.StdEncoding.EncodeToString([]byte(id_secret))

  r.Header.Add("content-type", "application/x-www-form-urlencoded")
  r.Header.Add("Authorization", auth_val)

  client := &http.Client{}
  res, err := client.Do(r)
  if err != nil {
    log.Fatal(err)
  }

  defer res.Body.Close()
  
  if res.StatusCode != http.StatusOK {
    errPost := &ErrorResponse{}
    if e := json.NewDecoder(res.Body).Decode(errPost); e != nil {
      log.Fatal(e)
    }

    log.Error("Authorization Error", "error", errPost.Error, "error description", errPost.ErrorDescription)
    return
  }

  access := &AccessResponse{}
  if e := json.NewDecoder(res.Body).Decode(access); e != nil {
    log.Fatal(e)
  }
  os.Setenv("SPOTIFY_ACCESS_CODE", access.AccessToken)
  os.Setenv("SPOTIFY_REFRESH_CODE", access.RefreshToken)
  
  log.Debug("Check code", "real", access.AccessToken, "env", os.Getenv("SPOTIFY_ACCESS_CODE"))
  log.Info("Updated access_code and refresh_token to env")
}

func (h *Handlers) HandleGetState(w http.ResponseWriter, r *http.Request) {
  u := &url.URL{
    Scheme: "https",
    Host: "api.spotify.com",
    Path: "v1/me/player",
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
    errPost := &ApiError{}
    if e := json.NewDecoder(res.Body).Decode(errPost); e != nil {
      log.Fatal(e)
    }

    log.Error("Error Occured", "error", errPost.Error)
    return
  }
  
  state := &PlayerState{}
  e := json.NewDecoder(res.Body).Decode(state)
  if e != nil {
    log.Fatal(e)
  }
  
  log.Info("Player State", "state", state)
  e = json.NewEncoder(w).Encode(state)
  if e != nil {
    log.Fatal(e)
  }
}




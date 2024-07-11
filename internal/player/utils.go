package player

import (
  "bytes"
  "encoding/base64"
  "encoding/json"
  "errors"
  "net/http"
  "net/url"
  "os"

  log "github.com/charmbracelet/log"
  yaml "gopkg.in/yaml.v3"
)

//Initialize Server
func StartServer() {
  log.SetLevel(log.DebugLevel)
  log.Info("Run TUI")

  //Get access codes
  f, err := os.ReadFile(".env.yaml")
  if err != nil {
    log.Fatal(err)
  }

  m := SpotifyKeys{}
  err = yaml.Unmarshal(f, &m)
  if err != nil {
    log.Fatal(err)
  }

  os.Setenv("SPOTIFY_CLIENT_ID", m.ClientID)
  os.Setenv("SPOTIFY_CLIENT_SECRET", m.ClientSecret)
  os.Setenv("SPOTIFY_REFRESH_TOKEN", m.RefreshToken)

  router := http.NewServeMux()
  LoadRoutes(router)

  server := http.Server{
    Addr: ":6969",
    Handler: router,
  }

  log.Info("Listening on port 6969")
  server.ListenAndServe()
}

//Request Functions
func (a *AccessRequest) MakeRequest() (AccessResponse, error) {
  u := &url.URL{
    Scheme: "https",
    Host: "accounts.spotify.com",
    Path: "api/token",
  }

  access := &AccessResponse{}
  body := make(url.Values)

  switch a.GrantType {
    case "authorization_code":     
    log.Info("Request by Code")
    body.Add("grant_type", a.GrantType)
    body.Add("code", a.AuthCode)
    body.Add("redirect_uri", REDIRECT_URI)
  case "refresh_token":
    log.Info("Request by Refresh Code")
    body.Add("grant_type", a.GrantType)
    body.Add("refresh_token", os.Getenv("SPOTIFY_REFRESH_TOKEN"))
    body.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
  default:
    return *access, errors.New("Error Reqesting Access Token.")
  }

  r, err := http.NewRequest("POST", u.String(), bytes.NewBuffer([]byte(body.Encode())))
  if err != nil {
    log.Fatal(err)
  }

  id_secret := string(os.Getenv("SPOTIFY_CLIENT_ID")+ ":" + os.Getenv("SPOTIFY_CLIENT_SECRET"))
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
    return *access, errors.New(errPost.Error)
  }

  if e := json.NewDecoder(res.Body).Decode(access); e != nil {
    log.Fatal(e)
  }

  return *access, nil
}

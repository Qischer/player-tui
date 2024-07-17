package player

//ID and Secret
type SpotifyKeys struct {
  ClientID      string    `yaml:"client_id"`
  ClientSecret  string    `yaml:"client_secret"`
  RefreshToken  string    `yaml:"refresh_token"`
}

//ERROR types 
type ErrorResponse struct {
  Error             string  `json:"error"`
  ErrorDescription  string  `json:"error_description"`
}

type ApiError struct {
  Error struct{
    Status  int64  `json:"status"`
    Message string `json:"message"`
  } `json:"error"`
}

// Requests | Responses
type AccessResponse struct {
  AccessToken    string `json:"access_token"`
  TokenType      string `json:"token_type"`
  Scope          string `json:"scope"` 
  ExpiresIn      int64  `json:"expires_in"`
  RefreshToken   string `json:"refresh_token"`
}

type AccessRequest struct {
  GrantType     string
  AuthCode      string
  RefreshToken  string
}

//Content
type PlayerState struct {
  ProgressMS int64        `json:"progress_ms"`
  Item       TrackObject  `json:"item"`
}

type TrackObject struct {
  Name      string    `json:"name"`
  Duration  int64     `json:"duration_ms"`
  Artists   []Artist  `json:"artists"`
}

type Artist struct {
  Name      string    `json:"name"`
}


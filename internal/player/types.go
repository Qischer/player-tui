package player

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

//Content
type PlayerState struct {
  ProgressMS int64        `json:"progress_ms"`
  Item       TrackObject  `json:"item"`
}

type TrackObject struct {
  Name      string    `json:"name"`
  Duration  int64     `json:"duration_ms"`
}

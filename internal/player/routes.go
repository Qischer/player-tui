package player

import "net/http"

func LoadRoutes(router *http.ServeMux) {
  h := &Handlers{}
  //router.HandleFunc("GET /", h.HandleIndex)
  router.HandleFunc("GET /auth", h.HandleAuth)
  router.HandleFunc("GET /callback", h.HandleCallback)
  router.HandleFunc("GET /player/state", h.HandleGetState)
  router.HandleFunc("PUT /player/play", h.HandlePlay)
  router.HandleFunc("PUT /player/pause", h.HandlePause)
}

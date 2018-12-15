package handle

import "net/http"

// CharCountHandler ...
type CharCountHandler struct {
	KitHandler http.Handler
}

func (h *CharCountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.KitHandler.ServeHTTP(w, r)
}

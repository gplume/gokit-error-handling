package handle

import "net/http"

// UpperCaseHandler handle http request for '/uppercase/**' routes
type UpperCaseHandler struct {
	KitHandler http.Handler
}

func (h *UpperCaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.KitHandler.ServeHTTP(w, r)
}

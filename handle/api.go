package handle

import (
	"fmt"
	"net/http"

	"github.com/gplume/no-mux/utils"
)

// Handlers ...
type Handlers struct {
	HomeHandler      http.Handler
	UpperCaseHandler http.Handler
	CharCountHandler http.Handler
}

func (h *Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.CutPath(r.URL.Path)
	switch head {
	case "":
		h.HomeHandler.ServeHTTP(w, r)
		return
	case "uppercase":
		h.UpperCaseHandler.ServeHTTP(w, r)
		return
	case "count":
		h.CharCountHandler.ServeHTTP(w, r)
		return
	}
	http.Error(w, fmt.Sprintf("Path: %q Not Found", r.URL.Path), http.StatusNotFound)
}

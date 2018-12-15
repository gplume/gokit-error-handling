package handle

import (
	"net/http"

	"github.com/gplume/no-mux/utils"
)

// Home ...
type Home struct{}

// HomeHandler implements ServeHTPP to return an http.Handler (interface satisfying)
func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, utils.JSMAP{
		"msg": "WELCOME TO THE API!",
	})
}

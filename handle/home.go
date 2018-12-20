package handle

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gplume/gokit-error-handling/utils"
)

// Home ...
type Home struct{}

// HomeHandler implements ServeHTPP to return an http.Handler (interface satisfying)
func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	head, _ := utils.CutPath(r.URL.Path)
	log.Println("head:", head)
	switch {
	case head == "":
		utils.JSON(w, http.StatusOK, utils.JSMAP{
			"msg": "WELCOME TO THE API!",
		})
		return
	default:
		switch r.URL.Query().Get(":ppat") {
		case "first":
			utils.JSON(w, http.StatusOK, utils.JSMAP{
				"msg": "WELCOME FIRST",
			})
			return
		case "second":
			utils.JSON(w, http.StatusOK, utils.JSMAP{
				"msg": "WELCOME SECOND",
			})
			return
		}
	}
	utils.JSON(w, http.StatusNotFound, utils.JSMAP{
		"msg": fmt.Sprintf("Route '%s' not found, so sorry, but BTW: WELCOME TO THE API!", r.URL.Path),
	})
	return
}

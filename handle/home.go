package handle

import (
	"fmt"
	"net/http"

	"github.com/gplume/gokit-error-handling/utils"
)

// Home ...
type Home struct{}

// HomeHandler implements ServeHTPP to return an http.Handler (interface satisfying)
func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	head, _ := utils.CutPath(r.URL.Path)
	fmt.Println("head:", head)
	if r.Method != http.MethodGet {
		utils.JSON(w, http.StatusNotFound, utils.JSMAP{
			"msg": fmt.Sprintf("Method %q, but BTW: WELCOME TO THE API!", r.Method),
		})
		return
	}
	switch {
	case head == "":
		utils.JSON(w, http.StatusOK, utils.JSMAP{
			"msg": "WELCOME TO THE API!",
		})
		return
	case head == "first":
		utils.JSON(w, http.StatusOK, utils.JSMAP{
			"msg": "WELCOME FIRST",
		})
		return
	case head == "second":
		utils.JSON(w, http.StatusOK, utils.JSMAP{
			"msg": "WELCOME SECOND",
		})
		return
	}
	utils.JSON(w, http.StatusNotFound, utils.JSMAP{
		"msg": fmt.Sprintf("Route '%s' not found, so sorry, but BTW: WELCOME TO THE API!", r.URL.Path),
	})
	return
}

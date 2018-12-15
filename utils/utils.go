package utils

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gplume/gokit-error-handling/stdlogger"
)

// JSMAP shortcut for map[string]interface{}
type JSMAP map[string]interface{}

// JSON ...
func JSON(w http.ResponseWriter, status int, value interface{}) {
	body, err := json.Marshal(value)
	if err != nil {
		stdlogger.Log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(status)
	_, err = w.Write(body)
}

// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
// https://husobee.github.io/golang/url-router/2015/06/15/why-do-all-golang-url-routers-suck.html

// CutPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func CutPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

package middle

import (
	"net/http"
)

// Wrapper ...
type Wrapper func(http.Handler) http.Handler

// Ware ...
func Ware(h http.Handler, wrappers ...Wrapper) http.Handler {
	// reverse order (onion):
	for _, warp := range wrappers {
		h = warp(h)
	}
	// or straight order:
	// for i := len(wrappers) - 1; i >= 0; i-- {
	// 	h = wrappers[i](h)
	// }
	return h
}

package middle

import (
	"fmt"
	"net/http"
	"time"
)

// Notify just log begining and ending of the called route (useless...)...
func Notify() Wrapper {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("––––––––––––––––––––– Start Route:", fmt.Sprintf("       %s        ", r.URL.Path), "––––––––––––––––––––––")
			defer func(begin time.Time) {
				fmt.Println("===================== Ends Route:   ", r.URL.Path, "after", time.Since(begin).Round(time.Nanosecond), "=====================")
			}(time.Now())
			h.ServeHTTP(w, r)
		})
	}
}

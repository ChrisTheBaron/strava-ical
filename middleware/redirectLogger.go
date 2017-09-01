package middleware

import (
	"github.com/golang/glog"
	"net/http"
)

//RedirectLogger logs if this request came from somewhere else.
type RedirectLogger func(http.Handler) http.HandlerFunc

func NewRedirectLogger() RedirectLogger {
	return RedirectLogger(func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Referer() != "" {
				glog.Infof("Redirected from %s", r.URL)
			}
			next.ServeHTTP(w, r)
		})
	})
}

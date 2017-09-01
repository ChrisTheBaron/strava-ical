package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/ChrisTheBaron/strava-ical/model"
	"github.com/ChrisTheBaron/strava-ical/utils"
	"github.com/golang/glog"
	"net/http"
)

const UserIdCtxKey = "UserIdCtxKey"

type VerifyJWT func(http.Handler) http.HandlerFunc

// VerifyJWT
func NewVerifyJWT(um *model.User, db *sql.DB, config *entities.Config) VerifyJWT {

	return VerifyJWT(func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// redirect accepts a slug ("login","calendar", etc.) and redirects to an absolute url
			// using the rootUrl and protocol in config.
			redirect := func(url string) {
				w.Header().Set("Location", fmt.Sprintf("%s://%s/%s", config.Protocol, config.RootUrl, url))
				w.WriteHeader(http.StatusTemporaryRedirect)
			}

			token := utils.GetTokenFromRequest(config, r)

			if token == "" {
				glog.Warningln("No Token set")
				redirect("/")
				return
			}

			uid, err := utils.ParseJWT(config, token)

			if err != nil {
				glog.Warningf("Failed to parse JWT: %s", err.Error())
				redirect("/")
				return
			}

			ok, err := um.ValidateUserId(uid)

			if err != nil {
				glog.Warningf("Failed to validate user: %s", err.Error())
				redirect("/")
				return
			}

			if ok {
				ctx := context.WithValue(r.Context(), UserIdCtxKey, uid)
				glog.Infof("Set context value: %d with key: %s", uid, UserIdCtxKey)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				glog.Warningln("Invalid user")
				redirect("/")
			}
		})
	})
}

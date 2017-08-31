package middleware

import (
	"context"
	"database/sql"
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

			redirect := func() {
				w.Header().Set("Location", config.Slugs.Login)
				w.WriteHeader(http.StatusForbidden)
			}

			token := getTokenFromRequest(config, r)

			if token == "" {
				glog.Warningln("No Token set")
				redirect()
				return
			}

			uid, err := utils.ParseJWT(config, token)

			if err != nil {
				glog.Warningf("Failed to parse JWT: %s", err.Error())
				redirect()
				return
			}

			ok, err := um.ValidateUserId(uid)

			if err != nil {
				glog.Warningf("Failed to validate user: %s", err.Error())
				redirect()
				return
			}

			if ok {
				ctx := context.WithValue(r.Context(), UserIdCtxKey, uid)
				glog.Infof("Set context value: %d with key: %s", uid, UserIdCtxKey)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				glog.Warningln("Invalid user")
				redirect()
			}
		})
	})
}

func getTokenFromRequest(c *entities.Config, r *http.Request) string {

	var token string

	token = getTokenFromHeader(c, r)

	if token == "" {
		token = getTokenFromCooke(c, r)
	}

	return token

}

func getTokenFromCooke(c *entities.Config, r *http.Request) string {

	cookie, err := r.Cookie(c.JWTCookieName)

	if err != nil {
		glog.Warning(err)
		return ""
	}

	return cookie.Value
}

func getTokenFromHeader(c *entities.Config, r *http.Request) string {

	return r.Header.Get("Bearer")

}

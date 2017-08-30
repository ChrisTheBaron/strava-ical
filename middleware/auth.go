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

// VerifyJWT
func VerifyJWT(um *model.User, db *sql.DB, config *entities.Config, next http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		redirect := func() {
			w.Header().Set("Location", config.Slugs.Login)
			w.WriteHeader(http.StatusForbidden)
		}

		cookie, err := r.Cookie(config.JWTCookieName)

		if err != nil {
			glog.Warningf("No Cookie set for %s", config.JWTCookieName)
			redirect()
			return
		}

		uid, err := utils.ParseJWT(config, cookie.Value)

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
}

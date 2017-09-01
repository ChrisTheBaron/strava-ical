package utils

import (
	"errors"
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/golang/glog"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"net/http"
)

func GenerateJWT(config *entities.Config, user entities.User) (string, error) {

	glog.Infof("Generating JWT for user id: %d, with key: %s", user.StravaId, config.JWTKey)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.StravaId,
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(config.JWTKey))

}

func ParseJWT(config *entities.Config, token string) (int, error) {

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.JWTKey), nil

	})

	if err != nil {
		return 0, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return int(claims["uid"].(float64)), nil
	} else {
		return 0, errors.New("Invalid JWT")
	}

}

func GetTokenFromRequest(c *entities.Config, r *http.Request) string {

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

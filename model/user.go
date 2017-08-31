package model

import (
	"database/sql"
	"github.com/ChrisTheBaron/strava-ical/entities"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	model
	db *sql.DB
}

// NewUser returns a new models.User
func NewUser(con *entities.Config, db *sql.DB) *User {
	return &User{model{config: con}, db}
}

func (m *User) ValidateUserId(id int) (bool, error) {
	result := m.db.QueryRow("SELECT count(1) FROM users WHERE strava_id = ?", id)
	var count int
	err := result.Scan(&count)
	return count > 0, err
}

func (m *User) Upsert(user entities.User) error {
	_, err := m.db.Exec("INSERT OR REPLACE INTO users (strava_id, firstname, lastname, email, strava_access_token) VALUES (?, ?, ?, ?, ?);",
		user.GetStravaId(), user.GetFirstname(), user.GetLastname(), user.GetEmail(), user.GetStravaAccessToken())
	return err
}

func (m *User) GetById(id int) (user entities.User, err error) {

	result := m.db.QueryRow("SELECT firstname, lastname, email, strava_id, strava_access_token FROM users WHERE strava_id = ?", id)

	var firstname string
	var lastname string
	var email string
	var strava_id int64
	var strava_access_token string

	err = result.Scan(&firstname, &lastname, &email, &strava_id, &strava_access_token)

	if err != nil {
		return
	}

	user = entities.NewUser(firstname, lastname, email, strava_id, strava_access_token)

	return

}

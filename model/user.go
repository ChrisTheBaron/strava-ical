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

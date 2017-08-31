package model

import (
	"database/sql"
	"github.com/ChrisTheBaron/strava-ical/entities"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
)

type Calendar struct {
	model
	db *sql.DB
}

// NewCalendar returns a new model.Calendars
func NewCalendar(con *entities.Config, db *sql.DB) *Calendar {
	return &Calendar{model{config: con}, db}
}

// Insert makes a new
func (m *Calendar) Insert(calendar entities.Calendar) (insertId uuid.UUID, err error) {

	if uuid.Equal(calendar.GetId(), uuid.Nil) {
		insertId = uuid.NewV4()
	} else {
		insertId = calendar.GetId()
	}

	_, err = m.db.Exec("INSERT INTO calendars (id, user_id) VALUES (?, ?)", insertId, calendar.GetUserId())

	return

}

func (m *Calendar) GetById(id uuid.UUID) (cal entities.Calendar, err error) {

	res := m.db.QueryRow("SELECT user_id FROM calendars WHERE id = ?", id.String())

	var uid int

	err = res.Scan(&uid)

	if err != nil {
		return
	}

	cal = entities.NewCalendar(id, uid)

	return

}

func (m *Calendar) GetAllForUser(uid int) (cals []entities.Calendar, err error) {

	res, err := m.db.Query("SELECT id, user_id FROM calendars WHERE user_id = ?", uid)

	if err != nil {
		return
	}

	for res.Next() {
		var user_id int
		var id uuid.UUID
		res.Scan(&id)
		res.Scan(&user_id)
		cals = append(cals, entities.NewCalendar(id, user_id))
	}

	return

}

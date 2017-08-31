package entities

import "github.com/satori/go.uuid"

type Calendar struct {
	id      uuid.UUID
	user_id int
}

func NewCalendar(id uuid.UUID, user_id int) Calendar {
	return Calendar{id, user_id}
}

func (c *Calendar) GetId() uuid.UUID {
	return c.id
}

func (c *Calendar) GetUserId() int {
	return c.user_id
}

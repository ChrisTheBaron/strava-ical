package entities

import "github.com/satori/go.uuid"

type Calendar struct {
	Id     uuid.UUID
	UserId int
}

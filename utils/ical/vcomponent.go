package ical

import "io"

type vComponent interface {
	Write(w io.Writer) error
}

// Package ical adapted from https://github.com/soh335/ical
package ical

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"time"
)

const (
	stampLayout    = "20060102T150405Z"
	dateLayout     = "20060102"
	dateTimeLayout = "20060102T150405"
)

type VCalendar struct {
	VERSION string // 2.0
	PRODID  string // -//My Company//NONSGML Event Calendar//EN
	URL     string // http://my.calendar/url

	NAME         string // My Calendar Name
	X_WR_CALNAME string // My Calendar Name
	DESCRIPTION  string // A description of my calendar
	X_WR_CALDESC string // A description of my calendar

	TIMEZONE_ID   string // Europe/London
	X_WR_TIMEZONE string // Europe/London

	REFRESH_INTERVAL string // PT12H
	X_PUBLISHED_TTL  string // PT12H

	CALSCALE string // GREGORIAN
	METHOD   string // PUBLISH

	VComponent []VComponent
}

func NewBasicVCalendar() *VCalendar {
	return &VCalendar{
		VERSION:  "2.0",
		CALSCALE: "GREGORIAN",
	}
}

func (c *VCalendar) AddComponent(comp VComponent) {
	c.VComponent = append(c.VComponent, comp)
}

func (c *VCalendar) Encode(w io.Writer) error {
	var b = bufio.NewWriter(w)

	if _, err := b.WriteString("BEGIN:VCALENDAR\r\n"); err != nil {
		return err
	}

	// use a slice map to preserve order during for range
	attrs := []map[string]string{
		{"VERSION:": c.VERSION},
		{"PRODID:": c.PRODID},
		{"URL:": c.URL},
		{"NAME:": c.NAME},
		{"X-WR-CALNAME:": c.X_WR_CALNAME},
		{"DESCRIPTION:": c.DESCRIPTION},
		{"X-WR-CALDESC:": c.X_WR_CALDESC},
		{"TIMEZONE-ID:": c.TIMEZONE_ID},
		{"X-WR-TIMEZONE:": c.X_WR_TIMEZONE},
		{"REFRESH-INTERVAL;VALUE=DURATION:": c.REFRESH_INTERVAL},
		{"X-PUBLISHED-TTL:": c.X_PUBLISHED_TTL},
		{"CALSCALE:": c.CALSCALE},
		{"METHOD:": c.METHOD},
	}

	for _, item := range attrs {
		for k, v := range item {
			if len(v) == 0 {
				continue
			}
			if _, err := b.WriteString(k + escapeCharacters(v) + "\r\n"); err != nil {
				return err
			}
		}
	}

	for _, component := range c.VComponent {
		if err := component.EncodeIcal(b); err != nil {
			return err
		}
	}

	if _, err := b.WriteString("END:VCALENDAR\r\n"); err != nil {
		return err
	}

	return b.Flush()
}

type VComponent interface {
	EncodeIcal(w io.Writer) error
}

type VEvent struct {
	VComponent
	UID         string
	DTSTAMP     time.Time
	DTSTART     time.Time
	DTEND       time.Time
	SUMMARY     string
	DESCRIPTION string
	TZID        string
	LOCATION    string

	AllDay bool
}

func (e VEvent) EncodeIcal(w io.Writer) error {

	var timeStampLayout, timeStampType, tzidTxt string

	if e.AllDay {
		timeStampLayout = dateLayout
		timeStampType = "DATE"
	} else {
		timeStampLayout = dateTimeLayout
		timeStampType = "DATE-TIME"
		if len(e.TZID) == 0 || e.TZID == "UTC" {
			timeStampLayout = timeStampLayout + "Z"
		}
	}

	if len(e.TZID) != 0 && e.TZID != "UTC" {
		tzidTxt = "TZID=" + e.TZID + ";"
	}

	b := bufio.NewWriter(w)
	if _, err := b.WriteString("BEGIN:VEVENT\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("DTSTAMP:" + e.DTSTAMP.UTC().Format(stampLayout) + "\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("UID:" + e.UID + "\r\n"); err != nil {
		return err
	}

	if len(e.TZID) != 0 && e.TZID != "UTC" {
		if _, err := b.WriteString("TZID:" + e.TZID + "\r\n"); err != nil {
			return err
		}
	}

	if _, err := b.WriteString("SUMMARY:" + escapeCharacters(e.SUMMARY) + "\r\n"); err != nil {
		return err
	}
	if e.DESCRIPTION != "" {
		if _, err := b.WriteString("DESCRIPTION:" + escapeCharacters(e.DESCRIPTION) + "\r\n"); err != nil {
			return err
		}
	}
	if e.LOCATION != "" {
		if _, err := b.WriteString("LOCATION:" + escapeCharacters(e.LOCATION) + "\r\n"); err != nil {
			return err
		}
	}
	if _, err := b.WriteString("DTSTART;" + tzidTxt + "VALUE=" + timeStampType + ":" + e.DTSTART.Format(timeStampLayout) + "\r\n"); err != nil {
		return err
	}

	if _, err := b.WriteString("DTEND;" + tzidTxt + "VALUE=" + timeStampType + ":" + e.DTEND.Format(timeStampLayout) + "\r\n"); err != nil {
		return err
	}

	if _, err := b.WriteString("END:VEVENT\r\n"); err != nil {
		return err
	}

	return b.Flush()
}

func escapeCharacters(s string) string {
	s = regexp.MustCompile(`(,|;|\\)`).ReplaceAllString(s, `\$1`)
	s = strings.Replace(s, "\n", `\n`, -1)
	return s
}

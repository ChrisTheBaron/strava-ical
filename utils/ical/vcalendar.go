// Package ical adapted from https://github.com/soh335/ical
package ical

import (
	"bufio"
	"github.com/golang/glog"
	"io"
	"regexp"
	"strings"
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

	vComponents []vComponent

	headersWritten bool
}

func (c *VCalendar) Add(comp vComponent) {
	c.vComponents = append(c.vComponents, comp)
}

func (c *VCalendar) WriteHeader(w io.Writer) error {

	if c.headersWritten {
		glog.Warningln("Headers already written")
		return nil
	}

	b := bufio.NewWriter(w)

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

	c.headersWritten = true

	return b.Flush()

}

func (c *VCalendar) Write(w io.Writer) error {

	if !c.headersWritten {
		glog.Warningln("Headers haven't been written")
		if err := c.WriteHeader(w); err != nil {
			return err
		}
	}

	b := bufio.NewWriter(w)

	for _, component := range c.vComponents {
		if err := component.Write(b); err != nil {
			return err
		}
	}

	if _, err := b.WriteString("END:VCALENDAR\r\n"); err != nil {
		return err
	}

	return b.Flush()
}

func escapeCharacters(s string) string {
	s = regexp.MustCompile(`(,|;|\\)`).ReplaceAllString(s, `\$1`)
	s = strings.Replace(s, "\n", `\n`, -1)
	return s
}

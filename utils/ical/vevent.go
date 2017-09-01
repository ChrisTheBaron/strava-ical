package ical

import (
	"bufio"
	"io"
	"time"
)

type VEvent struct {
	UID         string
	DTSTAMP     time.Time
	DTSTART     time.Time
	DTEND       time.Time
	SUMMARY     string
	DESCRIPTION string
	TZID        string
	LOCATION    string
	AllDay      bool
}

func (e VEvent) Write(w io.Writer) error {

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

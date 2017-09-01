package ical

import (
	"bufio"
	"fmt"
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
	if _, err := b.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", e.DTSTAMP.UTC().Format(stampLayout))); err != nil {
		return err
	}
	if _, err := b.WriteString(fmt.Sprintf("UID:%s\r\n", e.UID)); err != nil {
		return err
	}

	if len(e.TZID) != 0 && e.TZID != "UTC" {
		if _, err := b.WriteString(fmt.Sprintf("TZID:%s\r\n", e.TZID)); err != nil {
			return err
		}
	}

	if _, err := b.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", escapeCharacters(e.SUMMARY))); err != nil {
		return err
	}

	if e.DESCRIPTION != "" {
		if _, err := b.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", escapeCharacters(e.DESCRIPTION))); err != nil {
			return err
		}
	}
	if e.LOCATION != "" {
		if _, err := b.WriteString(fmt.Sprintf("LOCATION:%s\r\n", escapeCharacters(e.LOCATION))); err != nil {
			return err
		}
	}
	if _, err := b.WriteString(fmt.Sprintf("DTSTART;%sVALUE=%s:%s\r\n", tzidTxt, timeStampType, e.DTSTART.Format(timeStampLayout))); err != nil {
		return err
	}

	if _, err := b.WriteString(fmt.Sprintf("DTEND;%sVALUE=%s:%s\r\n", tzidTxt, timeStampType, e.DTEND.Format(timeStampLayout))); err != nil {
		return err
	}

	if _, err := b.WriteString("END:VEVENT\r\n"); err != nil {
		return err
	}

	return b.Flush()
}

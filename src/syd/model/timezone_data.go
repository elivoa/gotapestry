package model

import (
	"fmt"
	"time"
)

// About TimeZone

type TimeZoneInfo struct {
	Offset int    // javascript timezone Offset; retrive by new Date().getTimezoneOffset();
	Hour   int    // offset Hour
	Minute int    // offset minutes
	Prefix string // "+" or "-" sign before strict.
}

func NewTimeZoneInfo(offset int) *TimeZoneInfo {
	tzi := &TimeZoneInfo{Offset: offset}
	tzi.Hour = -offset / 60
	tzi.Minute = -offset % 60
	tzi.Prefix = "-"
	if tzi.Hour > 0 || tzi.Minute > 0 {
		tzi.Prefix = "+"
	}
	return tzi
}

func (m *TimeZoneInfo) String() string {
	hour := m.Hour
	if hour < 0 {
		hour = -hour
	}
	minute := m.Minute
	if minute < 0 {
		minute = -minute
	}
	return fmt.Sprintf("%s%d%d", m.Prefix, hour, minute)
}

func (m *TimeZoneInfo) Local(t time.Time) time.Time {
	return t.Add(time.Hour*time.Duration(m.Hour) + time.Minute*time.Duration(m.Minute))
}

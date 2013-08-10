package gxl

import (
	"fmt"
	"time"
)

func PrettyDay(then time.Time) string {
	y, m, d := then.Date()
	y2, m2, d2 := time.Now().Date()
	if y == y2 && m == m2 && d == d2 {
		return "今天"
	}
	then2 := time.Now().AddDate(0, 0, -1)
	y2, m2, d2 = then2.Date()
	if y == y2 && m == m2 && d == d2 {
		return "昨天"
	}
	return then.Format("2006-01-02")
}

// Seconds-based time units
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

// Time formats a time into a relative string.
// Time(someT) -> "3 weeks ago"
func HumanizeTime(then time.Time) string {
	now := time.Now()

	lbl := "ago"
	diff := now.Unix() - then.Unix()
	if then.After(now) {
		lbl = "from now"
		diff = then.Unix() - now.Unix()
	}

	switch {

	case diff <= 0:
		return "now"
	case diff <= 2:
		return fmt.Sprintf("1 second %s", lbl)
	case diff < 1*Minute:
		return fmt.Sprintf("%d seconds %s", diff, lbl)

	case diff < 2*Minute:
		return fmt.Sprintf("1 minute %s", lbl)
	case diff < 1*Hour:
		return fmt.Sprintf("%d minutes %s", diff/Minute, lbl)

	case diff < 2*Hour:
		return fmt.Sprintf("1 hour %s", lbl)
	case diff < 1*Day:
		return fmt.Sprintf("%d hours %s", diff/Hour, lbl)

	case diff < 2*Day:
		return fmt.Sprintf("1 day %s", lbl)
	case diff < 1*Week:
		return fmt.Sprintf("%d days %s", diff/Day, lbl)

	case diff < 2*Week:
		return fmt.Sprintf("1 week %s", lbl)
	case diff < 1*Month:
		return fmt.Sprintf("%d weeks %s", diff/Week, lbl)

	case diff < 2*Month:
		return fmt.Sprintf("1 month %s", lbl)
	case diff < 1*Year:
		return fmt.Sprintf("%d months %s", diff/Month, lbl)

	case diff < 18*Month:
		return fmt.Sprintf("1 year %s", lbl)
	}
	return then.String()
}

func HumanizeTimeCN(then time.Time) string {
	now := time.Now()

	lbl := "前"
	diff := now.Unix() - then.Unix()
	if then.After(now) {
		lbl = "以后"
		diff = then.Unix() - now.Unix()
	}

	switch {

	case diff <= 0:
		return "现在"
	case diff <= 2:
		return fmt.Sprintf("1 秒%s", lbl)
	case diff < 1*Minute:
		return fmt.Sprintf("%d 秒%s", diff, lbl)

	case diff < 2*Minute:
		return fmt.Sprintf("1 分钟%s", lbl)
	case diff < 1*Hour:
		return fmt.Sprintf("%d 分钟%s", diff/Minute, lbl)

	case diff < 2*Hour:
		return fmt.Sprintf("1 小时%s", lbl)
	case diff < 1*Day:
		return fmt.Sprintf("%d 小时%s", diff/Hour, lbl)

	case diff < 2*Day:
		return fmt.Sprintf("1 天%s", lbl)
	case diff < 1*Week:
		return fmt.Sprintf("%d 天%s", diff/Day, lbl)

	case diff < 2*Week:
		return fmt.Sprintf("1 周%s", lbl)
	case diff < 1*Month:
		return fmt.Sprintf("%d 周%s", diff/Week, lbl)

	case diff < 2*Month:
		return fmt.Sprintf("1 个月%s", lbl)
	case diff < 1*Year:
		return fmt.Sprintf("%d 个月%s", diff/Month, lbl)

	case diff < 18*Month:
		return fmt.Sprintf("1 年%s", lbl)
	}
	return then.String()
}

package ext

import "time"

const (
	// TimeFormatDefault default
	TimeFormatDefault = "2006-01-02 15:04:05"
	// TimeFormatyymmdd yymmdd
	TimeFormatyymmdd = "060102"
	// TimeFormatyyyymmdd yyyymmdd
	TimeFormatyyyymmdd = "20060102"
	// TimeFormatyyyymm yyyymm
	TimeFormatyyyymm = "200601"
)

//Today today 00:00:00
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// CurrHour current hour unix
func CurrHour() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
}

// CurrMinute minute
func CurrMinute() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
}

// CurrSecond second
func CurrSecond() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())
}

// CurrNanoSecond nanosecond
func CurrNanoSecond() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
}

// Formatyymmdd yymmdd
func Formatyymmdd(date time.Time) string {
	return date.Format(TimeFormatyymmdd)
}

// Formatyyyymmdd yyyymmdd
func Formatyyyymmdd(date time.Time) string {
	return date.Format(TimeFormatyyyymmdd)
}

// Formatyyyymm yyyymm
func Formatyyyymm(date time.Time) string {
	return date.Format(TimeFormatyyyymm)
}

// FormatDefault default
func FormatDefault(date time.Time) string {
	return date.Format(TimeFormatDefault)
}

//TicksToTime c# Ticks to time.Time
func TicksToTime(ticks int64) time.Time {
	ticks = ticks / 10
	n := int64(1000000)
	return time.Unix(ticks/n, ticks-(ticks/n)*n).AddDate(-1969, 0, 0).Add(-8 * time.Hour)
}

// TimeToTicks time to ticks
func TimeToTicks(t time.Time) int64 {
	return t.AddDate(1969, 0, 0).Unix() * 10000000
}

//TicksToUnixNano c# Ticks to UnixNano
func TicksToUnixNano(ticks int64) int64 {
	return TicksToTime(ticks).UnixNano()
}

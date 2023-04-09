package fn

import "time"

// IsSameDay 检查是否同一天
func IsSameDay(t1, t2 int64) bool {
	if ABSInt64(t1-t2) > 86400 {
		return false
	}

	date1 := time.Unix(t1, 0)
	date2 := time.Unix(t2, 0)
	_, _, d1 := date1.Date()
	_, _, d2 := date2.Date()

	f := false
	if d1 == d2 {
		f = (date1.Hour() >= DayStartHour && date2.Hour() >= DayStartHour) || (date1.Hour() < DayStartHour && date2.Hour() < DayStartHour)
	} else {
		f = (date1.Hour() >= DayStartHour && date2.Hour() < DayStartHour) || (date1.Hour() < DayStartHour && date2.Hour() >= DayStartHour)
	}

	return f
}

// ABSInt64 ...
func ABSInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// GetDay0Time 获取当天的0点时间
func GetDay0Time() int64 {
	d := time.Now()
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()).Unix()
}

// GetTimeStampDay0Time ...
func GetTimeStampDay0Time(tm int64) int64 {
	d := time.Unix(tm, 0)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()).Unix()
}

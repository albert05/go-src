package dates

import (
	"strconv"
	"time"
)

const DefaultDateFormatSTR = "2006-01-02 15:04:05"

func NowDateStr() string {
	return time.Now().Format(DefaultDateFormatSTR)
}

func NowDateShortStr() string {
	return time.Now().Format("20060102150405")
}

func NowTime() int64 {
	return time.Now().Unix()
}

func CurrentMicro() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

func TimeInt2float(t int64) float64 {
	r := t / 1e6
	micro := t - r*1e6

	tm := time.Unix(r, 0)
	f, err := strconv.ParseFloat(tm.Format("150405"), 64)
	if err != nil {
		return 0
	}

	return f + float64(micro)/1e6
}

func SleepSecond(t float64) {
	time.Sleep(time.Duration(t*1000) * time.Millisecond)
}

func NowYearMonthStr() string {
	return time.Now().Format("200601")
}

func NowDayStr() string {
	return time.Now().Format("20060102")
}

func TimeAfterNDays(n int) string {
	return time.Now().AddDate(0,0,n).Format(DefaultDateFormatSTR)
}

func DateAfterNDays(n int) string {
	return time.Now().AddDate(0,0,n).Format("2006-01-02") + " 00:00:00"
}

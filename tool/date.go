package tool

import (
	"time"
	"strconv"
)

const DEFAULT_DATE_FORMAT_STR = "2006-01-02 15:04:05"

func NowDateStr() string {
	return time.Now().Format(DEFAULT_DATE_FORMAT_STR)
}

func NowTime() int64 {
	return time.Now().Unix()
}

func CurrentMicro() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

func TimeInt2float(t int64) float64 {
	r := t / 1e6
	micro := t - r * 1e6

	tm := time.Unix(r, 0)
	f, err := strconv.ParseFloat(tm.Format("150405"), 64)
	if err != nil {
		return 0
	}

	return f + float64(micro) / 1e6
}

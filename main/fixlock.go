package main

import (
	"fmt"
	"kd.explorer/common"
	"kd.explorer/config"
	"kd.explorer/util/dates"
)

func main() {
	path := common.GetLockPath()
	files := common.GetAllFileByPattern(path, "RUN.MONITOR.TRANSFERS*")

	var fixTime int64 = config.RunDURATION + 30
	if len(files) > 0 {
		now := dates.NowTime()
		for _, file := range files {
			mt := common.GetFileModTime(file)
			if now-mt > fixTime {
				fmt.Println(common.RemoveFile(file))
			}
		}
	}
}

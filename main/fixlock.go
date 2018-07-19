package main

import (
	"kd.explorer/common"
	"kd.explorer/util/dates"
	"fmt"
	"kd.explorer/config"
)

func main() {
	path := common.GetLockPath()
	files := common.GetAllFileByPattern(path, "RUN.MONITOR.TRANSFERS*")

	var fixTime int64 = config.RunDURATION + 30
	if len(files) > 0 {
		now := dates.NowTime()
		for _, file := range files {
			mt := common.GetFileModTime(file)
			if now - mt > fixTime {
				fmt.Println(common.RemoveFile(file))
			}
		}
	}
}

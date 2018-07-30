package exception

import (
	"fmt"
	"kd.explorer/common"
	"kd.explorer/config"
	"kd.explorer/util/mail"
	"os"
	"runtime"
)

func Handle(isExit bool) {
	err := recover()
	if err != nil {
		buf := make([]byte, 2048)
		runtime.Stack(buf, true)
		errMsg := fmt.Sprintf("\n%s", buf)
		fmt.Println(err)
		fmt.Println(errMsg)

		var sErr string
		switch errObj := err.(type) {
		case error:
			sErr = errObj.Error()
		case string:
			sErr = errObj
		default:
			sErr = errMsg
		}

		mail.SendSingle(config.AdminMailer, config.ProNAME + "SYSTEM NOTICE", sErr)

		if isExit {
			common.UnLock()
			os.Exit(1)
		}
	}
}

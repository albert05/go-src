package tasks

const ShakeURL = "https://deposit.koudailc.com/daily-shake/daily-shake-award"
const ShareSuccessNUM = 3
const ShareTryNUM = 10

func Share(cookie string, params map[string]string) (bool, string) {
	succNum := 0
	for i := 1;;i++ {
		isOk, errMsg := Go(ShakeURL, cookie, params)
		if isOk {
			succNum++
		}

		if succNum >= ShareSuccessNUM || i >= ShareTryNUM {
			return isOk, errMsg
		}
	}
}

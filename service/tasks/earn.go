package tasks

//const EarnURL = "https://deposit.koudailc.com/user-level/earn"
const EarnURL = "https://deposit.koudailc.com/user-level/go-sign-in"
const EarnFailedRetryNUM = 3

func Earn(cookie string, params map[string]string) (bool, string) {
	for i := 1;;i++ {
		if isOk, errMsg := Go(EarnURL, cookie, params); isOk || i >= EarnFailedRetryNUM {
			return isOk, errMsg
		}
	}
}

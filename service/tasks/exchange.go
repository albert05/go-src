package tasks

const ExchangeURL = "https://deposit.koudailc.com/user-order-form/convert"

func Exchange(cookie string, params map[string]string) (bool, string) {
	return Go(ExchangeURL, cookie, params)
}

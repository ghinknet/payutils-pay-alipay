package payalipay

import "github.com/ghinknet/payutils/v3/driver"

func init() {
	driver.RegisterPay(Name, Driver{})
}

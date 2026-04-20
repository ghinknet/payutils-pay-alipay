package payalipay

import "github.com/ghinknet/payutils/v3/errors"

var ErrAlipayRespCodeInvalid = errors.New("alipay resp code invalid")
var ErrAlipayNotifyVerifyFailed = errors.New("alipay notify verify failed")

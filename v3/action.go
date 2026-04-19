package payalipay

import "github.com/ghinknet/payutils/v3/model"

func (c Client) Status(tradeID string) (model.ReturnStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (c Client) Close(tradeID string) error {
	//TODO implement me
	panic("implement me")
}

func (c Client) Refund(tradeID string, currency string, refundID string, totalAmount int64, refundAmount int64, reason string) error {
	//TODO implement me
	panic("implement me")
}

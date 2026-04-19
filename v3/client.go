package payalipay

import (
	"github.com/ghinknet/payutils/v3/model"
)

type Client struct{}
type Driver struct{}

func (d Driver) NewClient(params model.PayDriverClientParam) (model.PayClient, error) {
	return Client{}, nil
}

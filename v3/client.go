package payalipay

import (
	"github.com/ghinknet/payutils/v3/model"
	"github.com/ghinknet/toolbox/expr"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay/v3"
)

type Client struct {
	Params model.PayDriverClientParam
	Client *alipay.ClientV3
}
type Driver struct{}

func (d Driver) NewClient(params model.PayDriverClientParam) (model.PayClient, error) {
	// Create alipay client
	client, err := alipay.NewClientV3(
		params.Credential[AppID],
		params.Credential[AppCertPrivateKey],
		params.Credential[IsProd] == True,
	)
	if err != nil {
		return nil, err
	}

	// Set alipay cert
	if err = client.SetCert(
		[]byte(params.Credential[AppCert]),
		[]byte(params.Credential[RootCert]),
		[]byte(params.Credential[PublicCert]),
	); err != nil {
		return nil, err
	}

	// Debug switch
	client.DebugSwitch = expr.Ternary(params.Debug, gopay.DebugSwitch(gopay.DebugOn), gopay.DebugOff)

	return Client{
		Params: params,
		Client: client,
	}, nil
}

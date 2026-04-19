package payalipay

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/ghinknet/payutils/v3/errors"
	"github.com/ghinknet/payutils/v3/model"
	"github.com/ghinknet/payutils/v3/utils/currency"
	"github.com/go-pay/gopay"
)

func (c Client) Status(tradeID string) (model.ReturnStatus, error) {
	// Prepare params
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no",
		strings.Join(
			[]string{c.Params.TradeIDPrefix, tradeID, c.Params.TradeIDSuffix}, "",
		),
	).Set("query_options", []string{"send_pay_date"})

	// Check Alipay
	aliRsp, err := c.Client.TradeQuery(context.Background(), bm)
	if err != nil {
		return model.ReturnStatus{}, err
	}

	// Check return status
	if aliRsp.StatusCode != http.StatusOK && aliRsp.ErrResponse.Code != "ACQ.TRADE_NOT_EXIST" {
		return model.ReturnStatus{},
			ErrAlipayRespCodeInvalid.
				WithUpstreamName(Name).
				WithUpstreamCode(aliRsp.ErrResponse.Code).
				WithUpstreamMessage(aliRsp.ErrResponse.Message).
				WithUpstreamResponse(aliRsp)
	}
	if aliRsp.ErrResponse.Code == "ACQ.TRADE_NOT_EXIST" {
		return model.ReturnStatus{}, errors.ErrTradeNotExist.
			WithUpstreamName(Name).
			WithUpstreamCode(aliRsp.ErrResponse.Code).
			WithUpstreamMessage(aliRsp.ErrResponse.Message).
			WithUpstreamResponse(aliRsp)
	}

	// Success
	if aliRsp.StatusCode == 200 {
		return model.ReturnStatus{
			TradeStatus: MapState(aliRsp.TradeStatus),
			Upstream:    Name,
			Time:        FormatTime(aliRsp.SendPayDate),
		}, nil
	}

	return model.ReturnStatus{}, ErrAlipayRespCodeInvalid.
		WithUpstreamName(Name).
		WithUpstreamResponse(aliRsp)
}

func (c Client) Close(tradeID string) error {
	// Prepare params
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no",
		strings.Join(
			[]string{c.Params.TradeIDPrefix, tradeID, c.Params.TradeIDSuffix}, "",
		),
	)

	// Close order in Alipay
	aliRsp, err := c.Client.TradeClose(context.Background(), bm)
	if err != nil {
		return err
	}

	// Check return status
	if aliRsp.StatusCode != http.StatusOK && aliRsp.ErrResponse.Code != "ACQ.TRADE_NOT_EXIST" {
		return ErrAlipayRespCodeInvalid.
			WithUpstreamName(Name).
			WithUpstreamCode(aliRsp.ErrResponse.Code).
			WithUpstreamMessage(aliRsp.ErrResponse.Message).
			WithUpstreamResponse(aliRsp)
	}
	if aliRsp.ErrResponse.Code == "ACQ.TRADE_NOT_EXIST" {
		return errors.ErrTradeNotExist.
			WithUpstreamName(Name).
			WithUpstreamCode(aliRsp.ErrResponse.Code).
			WithUpstreamMessage(aliRsp.ErrResponse.Message).
			WithUpstreamResponse(aliRsp)
	}

	return nil
}

func (c Client) Refund(tradeID string, curr string, refundID string, totalAmount int64, refundAmount int64, reason string) error {
	// Check currency
	if curr != currency.CNY {
		return errors.ErrUnsupportedCurrency
	}

	// Prepare params
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no",
		strings.Join([]string{c.Params.TradeIDPrefix, tradeID, c.Params.TradeIDSuffix}, ""),
	)
	bm.Set("refund_amount", currency.CentsToYuan(refundAmount))
	bm.Set("out_request_no",
		strings.Join([]string{c.Params.TradeIDPrefix, tradeID, c.Params.TradeIDSuffix}, ""),
	)

	// Set reason
	if reason != "" {
		bm.Set("refund_reason", reason)
	}

	// Refund order in Alipay
	aliRsp, err := c.Client.TradeRefund(context.Background(), bm)
	if err != nil {
		return err
	}

	// Check return status
	if aliRsp.StatusCode != http.StatusOK && aliRsp.ErrResponse.Code != "ACQ.TRADE_NOT_EXIST" {
		return ErrAlipayRespCodeInvalid.
			WithUpstreamName(Name).
			WithUpstreamCode(aliRsp.ErrResponse.Code).
			WithUpstreamMessage(aliRsp.ErrResponse.Message).
			WithUpstreamResponse(aliRsp)
	}
	if aliRsp.ErrResponse.Code == "ACQ.TRADE_NOT_EXIST" {
		return errors.ErrTradeNotExist.
			WithUpstreamName(Name).
			WithUpstreamCode(aliRsp.ErrResponse.Code).
			WithUpstreamMessage(aliRsp.ErrResponse.Message).
			WithUpstreamResponse(aliRsp)
	}

	return nil
}

// MapState provides a method to map trade state (string) to internal trade state
func MapState(state string) model.TradeState {
	// Try to find in map
	internalState, ok := TradeStateMap[state]
	if !ok {
		return model.TradeStateUnknown
	}
	return internalState
}

// FormatTime provides a method to map trade time (string) to internal time
func FormatTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Time{}
	}

	// Timezone for alipay
	// That's suck, an international production
	// use UTC+8 as default timezone without any notice???
	// Even a monkey knows use standard time format
	// in an API return result
	// Fuck them all
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// Downgrade
		shanghai = time.FixedZone("CST", 8*60*60)
	}

	// Format time
	timeObj, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, shanghai)
	if err != nil {
		return time.Time{}
	}

	return timeObj
}

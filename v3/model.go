package payalipay

import "github.com/ghinknet/payutils/v3/model"

const (
	// TradeStateWaitBuyerPay Trade created and waiting for buyer to pay
	// -> TradeStatePending
	TradeStateWaitBuyerPay string = "WAIT_BUYER_PAY"
	// TradeStateSuccess Trade successes
	// -> TradeStateSuccess
	TradeStateSuccess string = "TRADE_SUCCESS"
	// TradeStateClosed Trade closed due to time out or refunded after pay
	// -> TradeStateClosed
	TradeStateClosed string = "TRADE_CLOSED"
	// TradeStateFinished Trade finished and cannot be refund
	// ->  TradeStateFinished
	TradeStateFinished string = "TRADE_FINISHED"
)

// TradeStateMap provides a map refer to internal trade state
var TradeStateMap = map[string]model.TradeState{
	TradeStateWaitBuyerPay: model.TradeStatePending,
	TradeStateSuccess:      model.TradeStateSuccess,
	TradeStateClosed:       model.TradeStateClosed,
	TradeStateFinished:     model.TradeStateFinished,
}

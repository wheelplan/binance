package binance

import (
	"fmt"
	"log"
	"math"
)

type PlacedOrder struct {
	Symbol              string `json:"symbol"`
	OrderId             int64  `json:"orderId"`
	ClientOrderId       string `json:"clientOrderId"`
	TransactTime        int64  `json:"transactTime"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeForce"`
	Type                string `json:"type"`
	IsIsolated          bool   `json:"isIsolated"` // 是否是逐仓 symbol 交易
	Side                string `json:"side"`
}

func (c *Client) MarginBuy(symbol string, quote float64) (res PlacedOrder, err error) {

	url := fmt.Sprintf("%s%s?symbol=%s&isIsolated=FALSE&side=BUY&sideEffectType=MARGIN_BUY&type=MARKET&quoteOrderQty=%f&newOrderRespType=RESULT", MarginBaseUrl, "/sapi/v1/margin/order", symbol, quote)

	_, err = c.do("POST", url, true, &res)
	if err != nil {
		return
	}

	return
}

func (c *Client) MarginSell(symbol string, quantity float64) (res PlacedOrder, err error) {

	url := fmt.Sprintf("%s%s?symbol=%s&isIsolated=FALSE&side=SELL&sideEffectType=AUTO_REPAY&type=MARKET&quantity=%f&newOrderRespType=RESULT", MarginBaseUrl, "/sapi/v1/margin/order", symbol, quantity)

	_, err = c.do("POST", url, true, &res)
	if err != nil {
		return
	}

	return
}

type CanceledOrder []struct {
	Status string `json:"status"`
}

func (c *Client) MarginCancelOrder(symbol string) (canceledOrder CanceledOrder, err error) {

	url := fmt.Sprintf("%s%s?symbol=%s&recvWindow=%d", MarginBaseUrl, "/sapi/v1/margin/openOrders", symbol, 30000)

	_, err = c.do("DELETE", url, true, &canceledOrder)
	if err != nil {
		return
	}

	return canceledOrder, nil
}

func (c *Client) MarginCancelOrderBuy(asset string) {

	_, err := c.MarginCancelOrder(asset + "USDT")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Canceled Order is OK .")
	}

	amount, err := c.MarginBalance("USDT")
	if err != nil {
		//messages.Wechat("CancelOrder Query Margin Account Balance Failed ")
		log.Println("CancelOrder Query Margin Account Balance req error", err)
		return
	}

	res, err := c.MarginBuy(asset+"USDT", math.Floor(amount))
	if err != nil {
		log.Println("MarginCancelOrderBuy req error", amount, err)
		//messages.Wechat("MarginCancelOrderBuy " + asset + " Failed ")
	} else {
		log.Printf("MarginCancelOrderBuy 花费 %v USDT, 成功买入 %v %s .\n%v", res.CummulativeQuoteQty, res.ExecutedQty, asset, res)
		//messages.Wechat("MarginCancelOrderBuy " + res.ExecutedQty + " " + asset + ", $" + res.CummulativeQuoteQty)
	}
}

func (c *Client) BorrowableBuy(asset string) (err error) {

	maxBorrowable, err := c.MaxBorrowable("USDT", "FALSE")
	if err != nil {
		//messages.Wechat("Query MaxBorrowable USDT Failed ")
		log.Println("Query MaxBorrowable req error", err)
		return
	}

	res, err := c.MarginBuy(asset+"USDT", math.Floor(maxBorrowable*0.99))
	if err != nil {
		//messages.Wechat("BorrowableBuy " + asset + " Failed .")
		log.Printf("BorrowableBuy %v req error.\n%v", maxBorrowable, err)
		return
	} else {
		log.Printf("BorrowableBuy 花费 %v USDT, 成功买入 %v %s .\n%v", res.CummulativeQuoteQty, res.ExecutedQty, asset, res)
		//messages.Wechat("BorrowableBuy " + res.ExecutedQty + " " + asset + ", $" + res.CummulativeQuoteQty)
	}

	return
}

func (c *Client) SuperMarginBuy(asset string) {

	// 查询杠杆账户 USDT 可用额度
	amount, err := c.MarginBalance("USDT")
	if err != nil {
		//messages.Wechat("Query Margin Account Balance Failed ")
		log.Println("Query Margin Account Balance req error", err)
	} else {
		log.Println("Margin USDT Balance is", amount)
	}

	// 杠杆账户全仓买入
	res, err := c.MarginBuy(asset+"USDT", math.Floor(amount))
	if err != nil {
		log.Println("MarginBuy req error", amount)
		//messages.Wechat("MarginBuy " + asset + " Failed ")
	} else {
		log.Printf("MarginBuy 花费 %v USDT, 成功买入 %v %s .\n%v", res.CummulativeQuoteQty, res.ExecutedQty, asset, res)
		//messages.Wechat("MarginBuy " + res.ExecutedQty + " " + asset + ", $" + res.CummulativeQuoteQty)
	}

	// 杠杆账户借贷后买入
	err = c.BorrowableBuy(asset)
	if err != nil {
		log.Println("BorrowableBuy req error", amount)
		//messages.Wechat("BorrowableBuy " + asset + " Failed ")
	}

	// 撤销订单后买入
	c.MarginCancelOrderBuy(asset)
}

func (c *Client) SuperMarginSell(asset string) (err error) {

	// 查询资产数量
	quantity, err := c.MarginBalance(asset)
	if err != nil {
		//messages.Wechat("Query Margin Account Balance Failed ")
		log.Println("Query Margin Account Balance req error", err)
		return
	}

	// 全仓卖出
	res, err := c.MarginSell(asset+"USDT", math.Floor(quantity))
	if err != nil {
		//messages.Wechat("MarginSell " + asset + " Failed ")
		log.Println("MarginSell req error", quantity, err)
	} else {
		log.Printf("MarginSell 成功卖出 %v %s, 收入 %v USDT .\n%v", res.ExecutedQty, asset, res.CummulativeQuoteQty, res)
		//messages.Wechat("MarginSell " + res.ExecutedQty + " " + asset + ", $" + res.CummulativeQuoteQty)
	}

	return
}

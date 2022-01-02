package binance

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"
	"wechat"
)

type CanceledContractOrder struct {
	Code uint16 `json:"code"`
	Msg  string `json:"msg"`
}

func (c *Client) CancelContractOrder(symbol string) (canceledOrder CanceledContractOrder, err error) {

	url := fmt.Sprintf("%s%s?symbol=%s&recvWindow=%d", ContractBaseUrl, "/fapi/v1/allOpenOrders", symbol, 30000)

	_, err = c.do("DELETE", url, true, &canceledOrder)
	if err != nil {
		return
	}

	return canceledOrder, nil
}

type Leveraged struct {
	Leverage         uint8  `json:"leverage"`         // 杠杆倍数
	MaxNotionalValue string `json:"maxNotionalValue"` // 当前杠杆倍数下允许的最大名义价值
	Symbol           string `json:"symbol"`           // 交易对
}

func (c *Client) Leverage(symbol string, leverage uint8) (res Leveraged, err error) {

	url := fmt.Sprintf("%s%s?symbol=%s&leverage=%d", ContractBaseUrl, "/fapi/v1/leverage", symbol, leverage)

	_, err = c.do("POST", url, true, &res)
	if err != nil {
		return
	}

	return
}

type ContractPlacedOrder struct {
	Symbol        string `json:"symbol"`        // 交易对
	OrderId       int64  `json:"orderId"`       // 系统订单号
	ClientOrderId string `json:"clientOrderId"` // 用户自定义的订单号
	Price         string `json:"price"`         // 委托价格
	OrigQty       string `json:"origQty"`       // 原始委托数量
	ExecutedQty   string `json:"executedQty"`   // 成交量
	Status        string `json:"status"`        // 订单状态
	TimeInForce   string `json:"timeForce"`     // 有效方法
	Type          string `json:"type"`          // 订单类型
	IsIsolated    bool   `json:"isIsolated"`    // 是否是逐仓 symbol 交易
	Side          string `json:"side"`          // 买卖方向
	CumQty        string `json:"cumQty"`
	CumQuote      string `json:"cumQuote"`      // 成交金额
	AvgPrice      string `json:"avgPrice"`      // 平均成交价
	ReduceOnly    bool   `json:"reduceOnly"`    // 仅减仓
	PositionSide  string `json:"positionSide"`  // 持仓方向
	StopPrice     string `json:"stopPrice"`     // 触发价，对`TRAILING_STOP_MARKET`无效
	ClosePosition bool   `json:"closePosition"` // 是否条件全平仓
	OrigType      string `json:"origType"`      // 触发前订单类型
	ActivatePrice string `json:"activatePrice"` // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate     string `json:"priceRate"`     // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	UpdateTime    int64  `json:"updateTime"`    // 更新时间
	WorkingType   string `json:"workingType"`   // 条件价格触发类型
	PriceProtect  bool   `json:"priceProtect"`  // 是否开启条件单触发保护
}

func (c *Client) ContractBuy(symbol string, quantity float64, position string) (res ContractPlacedOrder, err error) {

	url := fmt.Sprintf("%s%s?symbol=%s&side=BUY&type=MARKET&quantity=%f&positionSide=%s", ContractBaseUrl, "/fapi/v1/order", symbol, quantity, position)

	_, err = c.do("POST", url, true, &res)
	if err != nil {
		return
	}

	return
}

func (c *Client) ContractSell(symbol string, quantity float64, position string) (res ContractPlacedOrder, err error) {

	url := fmt.Sprintf("%s%s?symbol=%s&side=SELL&type=MARKET&quantity=%f&positionSide=%s", ContractBaseUrl, "/fapi/v1/order", symbol, quantity, position)

	_, err = c.do("POST", url, true, &res)
	if err != nil {
		return
	}

	return
}

func (c *Client) ContractStopSell(symbol, typeOf, position string, stopPrice float64) (res ContractPlacedOrder, err error) {

	url := fmt.Sprintf("%s%s?symbol=%s&side=SELL&type=%s&closePosition=true&positionSide=%s&stopPrice=%f", ContractBaseUrl, "/fapi/v1/order", symbol, typeOf, position, stopPrice)

	_, err = c.do("POST", url, true, &res)
	if err != nil {
		return
	}

	return
}

func (c *Client) SuperContractBuy(symbol string, quantity float64) {

	startTime := time.Now()

	// 全仓买入
	tasks := []func(){
		func() {
			_, err := c.ContractBuy(symbol, quantity*3, "LONG")
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Successfully purchased", quantity*3)
			}
		},
		func() {
			_, err := c.ContractBuy(symbol, quantity*2, "LONG")
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Successfully purchased", quantity*2)
			}
		},
		func() {
			_, err := c.ContractBuy(symbol, quantity*2, "LONG")
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Successfully purchased", quantity*2)
			}
		},
		func() {
			_, err := c.ContractBuy(symbol, quantity*1.5, "LONG")
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Successfully purchased", quantity*1.5)
			}
		},
		func() {
			_, err := c.ContractBuy(symbol, quantity*1.5, "LONG")
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Successfully purchased", quantity*1.5)
			}
		},
		func() {
			_, err := c.ContractBuy(symbol, quantity, "LONG")
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Successfully purchased", quantity)
			}
		},
		func() {
			_, err := c.ContractBuy(symbol, quantity, "LONG")
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Successfully purchased", quantity)
			}
		},
		func() {
			_, err := c.ContractBuy(symbol, quantity, "LONG")
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Successfully purchased", quantity)
			}
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for _, task := range tasks {
		task := task

		go func() {
			defer wg.Done()
			task()
		}()
	}

	wg.Wait()

	log.Println(time.Since(startTime))

	log.Println("lucky ...")

	account, err := c.AccountData()
	log.Println("total wallet balance is", account.TotalWalletBalance, "total unrealized profit is", account.TotalUnrealizedProfit)

	var positionAmt, entryPrice float64
	for _, position := range account.Positions {
		if position.Symbol == symbol && position.PositionSide == "LONG" {
			positionAmt = position.PositionAmt
			entryPrice = position.EntryPrice
			log.Println(position.PositionAmt, position.EntryPrice)
		}
	}

	// 止损
	r, err := c.ContractStopSell("DOGEUSDT", "STOP_MARKET", "LONG", entryPrice*0.985)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Stop loss price is", r.StopPrice)
	}

	// 止盈
	r, err = c.ContractStopSell("DOGEUSDT", "TAKE_PROFIT_MARKET", "LONG", entryPrice*1.09)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Take profit price is", r.StopPrice)
	}

	PositionAmt := fmt.Sprintf("%v", positionAmt)
	EntryPrice := fmt.Sprintf("%v", entryPrice)
	CumQuote := fmt.Sprintf("%v", positionAmt*entryPrice)

	err = wechat.BuyDogeCard(PositionAmt, EntryPrice, CumQuote)
	if err != nil {
		log.Println(err)
	}

}

func (c *Client) SuperContractSell(symbol string) (err error) {

	account, err := c.AccountData()
	if err != nil {
		log.Println(err)
	}
	log.Println("total wallet balance is", account.TotalWalletBalance, "total unrealized profit is", account.TotalUnrealizedProfit)

	var positionAmt float64
	for _, position := range account.Positions {
		if position.Symbol == symbol && position.PositionSide == "LONG" {
			positionAmt = position.PositionAmt
		}
	}

	res, err := c.ContractSell(symbol, math.Floor(positionAmt), "LONG")
	if err != nil {
		log.Println(res)
		return
	}

	account, err = c.AccountData()
	if err != nil {
		log.Println(err)
	}

	amount := fmt.Sprintf("%v", positionAmt)
	balance := fmt.Sprintf("%v", account.TotalWalletBalance)

	err = wechat.SellDogeCard(amount, balance)
	if err != nil {
		log.Println(err)
	}

	return
}

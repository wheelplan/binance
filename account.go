package binance

import (
	"fmt"
)

type UserData struct {
	TotalWalletBalance    float64 `json:"totalWalletBalance,string"`    // 账户总余额
	TotalUnrealizedProfit float64 `json:"totalUnrealizedProfit,string"` // 持仓未实现盈亏总额
	Positions             []struct {
		Symbol       string  `json:"symbol"`             // 持仓交易对
		Isolated     bool    `json:"isolated"`           // 是否是逐仓模式
		PositionSide string  `json:"positionSide"`       // 持仓方向
		EntryPrice   float64 `json:"entryPrice,string"`  // 持仓成本价
		PositionAmt  float64 `json:"positionAmt,string"` // 持仓数量
	} `json:"positions"`
}

func (c *Client) AccountData() (res UserData, err error) {

	url := fmt.Sprintf("%s%s", ContractBaseUrl, "/fapi/v2/account")

	_, err = c.do("GET", url, true, &res)
	if err != nil {
		return
	}

	return
}

package binance

import (
	"fmt"
)

type TickerPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

func (c *Client) MarginPrice(symbol string) (res float64, err error) {

	ticker := TickerPrice{}

	url := fmt.Sprintf("%s%s?symbol=%s", MarginBaseUrl, "/api/v3/ticker/price", symbol)

	_, err = c.do("GET", url, false, &ticker)
	if err != nil {
		return
	}

	return ticker.Price, nil
}

func (c *Client) ContractPrice(symbol string) (res float64, err error) {

	ticker := TickerPrice{}

	url := fmt.Sprintf("%s%s?symbol=%s", ContractBaseUrl, "/fapi/v1/ticker/price", symbol)

	_, err = c.do("GET", url, false, &ticker)
	if err != nil {
		return
	}

	return ticker.Price, nil
}

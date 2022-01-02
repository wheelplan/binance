package binance

import (
	"fmt"
)

type Borrowable struct {
	Amount float64 `json:"amount,string"`
}

func (c *Client) MaxBorrowable(asset, isolatedSymbol string) (res float64, err error) {

	borrowable := Borrowable{}

	url := fmt.Sprintf("%s%s?asset=%s&isIsolated=%s", MarginBaseUrl, "/sapi/v1/margin/maxBorrowable", asset, isolatedSymbol)

	_, err = c.do("GET", url, true, &borrowable)
	if err != nil {
		return
	}

	return borrowable.Amount, nil
}

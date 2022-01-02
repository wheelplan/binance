package market

import "fmt"

func (c *Client) SpotBuy(symbol string, quote float64) (res PlacedOrder, err error) {

    url := fmt.Sprintf("%s%s?symbol=%s&side=BUY&type=MARKET&quoteOrderQty=%f&newOrderRespType=RESULT", MarginBaseUrl, "/api/v3/order", symbol, quote)

    _, err = c.do("POST", url, true, &res)
    if err != nil {
        return
    }

    return
}

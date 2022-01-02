package binance

import (
	"errors"
	"fmt"
	"log"
)

type ID struct {
	TranId int64 `json:"tranId"`
}

func (c *Client) Transfer(transferType, asset string, amount float64) (tranId int64, err error) {

	url := fmt.Sprintf("%s%s?type=%s&asset=%s&amount=%f", MarginBaseUrl, "/sapi/v1/asset/transfer", transferType, asset, amount)

	id := ID{}

	_, err = c.do("POST", url, true, &id)
	if err != nil {
		return
	}

	return id.TranId, nil
}

func (c *Client) UMarginRansfer(asset string) (err error) {

	count := true

	for i := 0; i < 3; i++ {

		qty, err := c.ContractBalance(asset)
		if err != nil {
			log.Println("ContractBalance", err)
		}

		_, err = c.Transfer("UMFUTURE_MARGIN", asset, qty*0.98)
		if err != nil {
			log.Println(err)
			continue
		} else {
			msg := fmt.Sprintf("UMarginRansfer %f %s", qty*0.97, asset)
			log.Printf(msg)
			//messages.Wechat(msg)
			count = false
			break
		}
	}

	if count {
		err = errors.New("UMarginRansfer is failed . ")
	}

	return
}

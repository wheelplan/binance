package market

import (
	"log"
	"time"
)

func (c *Client) Allocate(asset string, waitTime int64, quantity float64) {

	startTime := time.Now()

	// 合约全仓买入
	c.SuperContractBuy("DOGEUSDT", quantity)

	log.Println(time.Since(startTime))

	// 等待时机
	time.Sleep(time.Duration(waitTime) * time.Second)

	// 合约账户平仓
	ContractSellStatus := c.SuperContractSell("DOGE")

	// Lucky
	if ContractSellStatus != nil {
		for n := 0; n < 99; n++ {
			//messages.Wechat("SellStatus is bad !")
			time.Sleep(time.Duration(6000) * time.Millisecond)
		}
	}
}

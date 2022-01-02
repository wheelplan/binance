package market

import (
	"fmt"
)

type marginAccount struct {
	TradeEnabled        bool    `json:"tradeEnabled"`
	TransferEnabled     bool    `json:"transferEnabled"`
	BorrowEnabled       bool    `json:"borrowEnabled"`
	MarginLevel         float64 `json:"marginLevel,string"`
	TotalAssetOfBtc     float64 `json:"totalAssetOfBtc,string"`
	TotalLiabilityOfBtc float64 `json:"totalLiabilityOfBtc,string"`
	TotalNetAssetOfBtc  float64 `json:"totalNetAssetOfBtc,string"`

	UserAssets []struct {
		Asset    string  `json:"asset"`
		Free     float64 `json:"free,string"`
		Locked   float64 `json:"locked,string"`
		Borrowed float64 `json:"borrowed,string"`
		Interest float64 `json:"interest,string"`
		NetAsset float64 `json:"netAsset,string"`
	} `json:"userAssets"`
}

func (c *Client) MarginBalance(asset string) (res float64, err error) {

	balance := marginAccount{}

	url := fmt.Sprintf("%s%s", MarginBaseUrl, "/sapi/v1/margin/account")

	_, err = c.do("GET", url, true, &balance)
	if err != nil {
		return
	}

	for i, _ := range balance.UserAssets {
		if balance.UserAssets[i].Asset == asset {
			return balance.UserAssets[i].Free, nil
		}
	}

	return
}

type ContractAccount struct {
	AvailableBalance float64 `json:"availableBalance,string"`
	Positions        []struct {
		Symbol       string  `json:"symbol"`
		PositionSide string  `json:"positionSide"`
		PositionAmt  float64 `json:"positionAmt,string"`
		Leverage     string  `json:"leverage"`
		Isolated     bool    `json:"isolated"`
	} `json:"positions"`
}

func (c *Client) ContractBalance(asset string) (res float64, err error) {

	balance := ContractAccount{}

	url := ContractBaseUrl + "/fapi/v2/account"

	_, err = c.do("GET", url, true, &balance)
	if err != nil {
		return
	}

	if asset == "USDT" {
		return balance.AvailableBalance, nil
	}

	for i, _ := range balance.Positions {
		if balance.Positions[i].Symbol == asset+"USDT" {
			return balance.Positions[i].PositionAmt, nil
		}
	}

	return
}

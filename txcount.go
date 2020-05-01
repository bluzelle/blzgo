package bluzelle

import (
	"encoding/json"
	"strconv"
)

func (ctx *Client) TxCount(gasInfo *GasInfo) (int, error) {
	transaction := &Transaction{
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/count",
		GasInfo:            gasInfo,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return 0, err
	}

	res := &CountResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return 0, err
	}

	if count, err := strconv.Atoi(res.Count); err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

package bluzelle

import (
	"encoding/json"
	"strconv"
)

func (ctx *Client) TxCount() (int, error) {
	transaction := &Transaction{
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/count",
		Client:             ctx,
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

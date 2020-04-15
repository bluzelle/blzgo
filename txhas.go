package bluzelle

import (
	"encoding/json"
)

func (ctx *Client) TxHas(key string) (bool, error) {
	transaction := &Transaction{
		Key:                key,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/has",
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return false, err
	}

	res := &HasResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return false, err
	}
	return res.Has, nil
}

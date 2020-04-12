package bluzelle

import (
	"encoding/json"
)

func (ctx *Client) TxRead(key string) (string, error) {
	transaction := &Transaction{
		Key:                key,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/read",
		Client:             ctx,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return "", err
	}

	res := &ReadResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return "", err
	}
	return res.Value, nil
}

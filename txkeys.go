package bluzelle

import (
	"encoding/json"
)

func (ctx *Client) TxKeys() ([]string, error) {
	transaction := &Transaction{
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/keys",
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return nil, err
	}

	res := &KeysResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

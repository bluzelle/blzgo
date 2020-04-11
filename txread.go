package bluzelle

import (
	"encoding/json"
)

func (ctx *Client) TxRead(key string, gasInfo *GasInfo) (string, error) {
	transaction := &Transaction{
		Key:                key,
		Address:            ctx.Options.Address,
		UUID:               ctx.Options.UUID,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/read",
		GasInfo:            gasInfo,
		ChainId:            ctx.Options.ChainId,
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

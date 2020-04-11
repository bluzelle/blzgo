package bluzelle

import (
	"encoding/json"
)

func (ctx *Client) TxHas(key string, gasInfo *GasInfo) (bool, error) {
	transaction := &Transaction{
		Key:                key,
		Address:            ctx.Options.Address,
		UUID:               ctx.Options.UUID,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/has",
		GasInfo:            gasInfo,
		ChainId:            ctx.Options.ChainId,
		Client:             ctx,
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

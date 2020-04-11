package bluzelle

import (
	"encoding/json"
)

func (ctx *Client) TxKeys(uuid string, gasInfo *GasInfo) ([]string, error) {
	transaction := &Transaction{
		Address:            ctx.Options.Address,
		UUID:               uuid,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/keys",
		GasInfo:            gasInfo,
		ChainId:            ctx.Options.ChainId,
		Client:             ctx,
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

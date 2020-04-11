package bluzelle

import (
	"encoding/json"
)

func (ctx *Client) TxKeyValues(uuid string, gasInfo *GasInfo) ([]*KeyValuesResponseResultKeyValue, error) {
	transaction := &Transaction{
		Address:            ctx.Options.Address,
		UUID:               uuid,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/keyvalues",
		GasInfo:            gasInfo,
		ChainId:            ctx.Options.ChainId,
		Client:             ctx,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return nil, err
	}

	res := &KeyValuesResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res.KeyValues, nil
}

package bluzelle

import (
	"encoding/json"
)

func (ctx *Client) TxKeyValues(gasInfo *GasInfo) ([]*KeyValue, error) {
	transaction := &Transaction{
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/keyvalues",
		GasInfo:            gasInfo,
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

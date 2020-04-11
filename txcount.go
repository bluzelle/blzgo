package bluzelle

import (
	"encoding/json"
	"strconv"
)

func (ctx *Client) TxCount(uuid string, gasInfo *GasInfo) (int, error) {
	transaction := &Transaction{
		Address:            ctx.Options.Address,
		UUID:               uuid,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/count",
		GasInfo:            gasInfo,
		ChainId:            ctx.Options.ChainId,
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

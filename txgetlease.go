package bluzelle

import (
	"encoding/json"
	"strconv"
)

func (ctx *Client) TxGetLease(key string, gasInfo *GasInfo) (int64, error) {
	transaction := &Transaction{
		Key:                key,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/getlease",
		GasInfo:            gasInfo,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return 0, err
	}

	res := &GetLeaseResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return 0, err
	}
	lease, err := strconv.Atoi(res.Lease)
	return int64(lease), err
}

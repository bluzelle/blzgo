package bluzelle

import (
	"encoding/json"
)

func (ctx *Client) TxGetNShortestLeases(n uint64) ([]*KeyLease, error) {
	transaction := &Transaction{
		N:                  n,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/getnshortestlease",
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return nil, err
	}

	res := &GetNShortestLeasesResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res.KeyLeases, nil
}

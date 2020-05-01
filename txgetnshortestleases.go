package bluzelle

import (
	"encoding/json"
)

func (ctx *Client) TxGetNShortestLeases(n uint64, gasInfo *GasInfo) ([]*GetNShortestLeasesResponseResultKeyLease, error) {
	transaction := &Transaction{
		N:                  n,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/getnshortestlease",
		GasInfo:            gasInfo,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return nil, err
	}

	res := &GetNShortestLeasesResponseResult{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	kl, err := res.GetHumanizedKeyLeases()
	return kl, err
}

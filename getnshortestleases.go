package bluzelle

import (
	"encoding/json"
	"fmt"
)

type GetNShortestLeasesResponseResult struct {
	KeyLeases []*KeyLease `json:"keyleases"`
}

type GetNShortestLeasesResponse struct {
	Result *GetNShortestLeasesResponseResult `json:"result"`
}

func (ctx *Client) GetNShortestLeases(n uint64) ([]*KeyLease, error) {
	body, err := ctx.APIQuery(fmt.Sprintf("/crud/getnshortestlease/%s/%d", ctx.options.UUID, n))
	if err != nil {
		return nil, err
	}

	res := &GetNShortestLeasesResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res.Result.KeyLeases, nil
}

package bluzelle

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type GetNShortestLeasesResponseResultKeyLease struct {
	Key   string
	Lease int64
}

type GetNShortestLeasesResponseResult struct {
	KeyLeases []*KeyLease `json:"keyleases"`
}

func (result *GetNShortestLeasesResponseResult) GetHumanizedKeyLeases() ([]*GetNShortestLeasesResponseResultKeyLease, error) {
	if len(result.KeyLeases) == 0 {
		return []*GetNShortestLeasesResponseResultKeyLease{}, nil
	}
	var ret []*GetNShortestLeasesResponseResultKeyLease
	for _, kl := range result.KeyLeases {
		if lease, err := strconv.ParseInt(kl.Lease, 10, 64); err != nil {
			return nil, err
		} else {
			ret = append(ret, &GetNShortestLeasesResponseResultKeyLease{Key: kl.Key, Lease: lease * BLOCK_TIME_IN_SECONDS})
		}
	}
	return ret, nil
}

type GetNShortestLeasesResponse struct {
	Result *GetNShortestLeasesResponseResult `json:"result"`
}

func (ctx *Client) GetNShortestLeases(n uint64) ([]*GetNShortestLeasesResponseResultKeyLease, error) {
	body, err := ctx.APIQuery(fmt.Sprintf("/crud/getnshortestlease/%s/%d", ctx.options.UUID, n))
	if err != nil {
		return nil, err
	}

	res := &GetNShortestLeasesResponse{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	kl, err := res.Result.GetHumanizedKeyLeases()
	return kl, err
}

package bluzelle

import (
	"encoding/json"
	"strconv"
)

type GetLeaseResponseResult struct {
	Lease string `json:"lease"`
}

type GetLeaseResponse struct {
	Result *GetLeaseResponseResult `json:"result"`
}

func (ctx *Client) GetLease(key string) (int64, error) {
	body, err := ctx.APIQuery("/crud/getlease/" + ctx.options.UUID + "/" + key)
	if err != nil {
		return 0, err
	}

	res := &GetLeaseResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return 0, err
	}

	lease, err := strconv.Atoi(res.Result.Lease)
	return int64(lease * BLOCK_TIME_IN_SECONDS), err
}

package bluzelle

import (
	"encoding/json"
)

type HasResponseResult struct {
	Has bool `json:"has"`
}

type HasResponse struct {
	Result *HasResponseResult `json:"result"`
}

func (ctx *Client) Has(key string) (bool, error) {
	body, err := ctx.APIQuery("/crud/has/" + ctx.options.UUID + "/" + key)
	if err != nil {
		return false, err
	}

	res := &HasResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return false, err
	}
	return res.Result.Has, nil
}

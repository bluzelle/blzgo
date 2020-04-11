package bluzelle

import (
	"encoding/json"
	"strconv"
)

type CountResponseResult struct {
	Count string `json:"count"`
}

type CountResponse struct {
	Result *CountResponseResult `json:"result"`
}

func (ctx *Client) Count() (int, error) {
	body, err := ctx.APIQuery("/crud/count/" + ctx.options.UUID)
	if err != nil {
		return 0, err
	}

	res := &CountResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return 0, err
	}

	if count, err := strconv.Atoi(res.Result.Count); err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

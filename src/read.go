package bluzelle

import (
	"encoding/json"
)

type ReadResponseResult struct {
	Value string `json:"value"`
	Key   string `json:"key"`
	UUID  string `json:"uuid"`
}

type ReadResponse struct {
	Result *ReadResponseResult `json:"result"`
}

func (ctx *Client) Read(key string, proove bool) (string, error) {
	res := &ReadResponse{}

	body, err := ctx.Query("/crud/read/" + ctx.UUID + "/" + key)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		return "", err
	}

	return res.Result.Value, nil
}

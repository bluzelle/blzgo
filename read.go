package bluzelle

import (
	"encoding/json"
)

type ReadResponseResult struct {
	Value string `json:"value"`
}

type ReadResponse struct {
	Result *ReadResponseResult `json:"result"`
}

func (ctx *Client) Read(key string) (string, error) {
	value, err := ctx.ReadWithProof(key, false)
	return value, err
}

func (ctx *Client) ReadWithProof(key string, prove bool) (string, error) {
	path := "read"
	if prove {
		path = "pread"
	}
	body, err := ctx.APIQuery("/crud/" + path + "/" + ctx.options.UUID + "/" + key)
	if err != nil {
		return "", err
	}

	res := &ReadResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return "", err
	}
	return res.Result.Value, nil
}

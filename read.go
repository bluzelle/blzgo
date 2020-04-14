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
	body, err := ctx.APIQuery("/crud/read/" + ctx.options.UUID + "/" + key)
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

func (ctx *Client) ProvenRead(key string) (string, error) {
	body, err := ctx.APIQuery("/crud/pread/" + ctx.options.UUID + "/" + key)
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

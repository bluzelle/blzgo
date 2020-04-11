package bluzelle

import (
	"encoding/json"
)

type KeysResponseResult struct {
	Keys []string `json:"keys"`
}

type KeysResponse struct {
	Result *KeysResponseResult `json:"result"`
}

func (ctx *Client) Keys(uuid string) ([]string, error) {
	body, err := ctx.APIQuery("/crud/keys/" + uuid)
	if err != nil {
		return nil, err
	}

	res := &KeysResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res.Result.Keys, nil
}
package bluzelle

import (
	"encoding/json"
)

type KeyValuesResponseResult struct {
	KeyValues []*KeyValue `json:"keyvalues"`
}

type KeyValuesResponse struct {
	Result *KeyValuesResponseResult `json:"result"`
}

func (ctx *Client) KeyValues() ([]*KeyValue, error) {
	body, err := ctx.APIQuery("/crud/keyvalues/" + ctx.options.UUID)
	if err != nil {
		return nil, err
	}

	res := &KeyValuesResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res.Result.KeyValues, nil
}

package bluzelle

import (
	"encoding/json"
)

type Account struct {
	PublicKey     string `json:"public_key"`
	AccountNumber int    `json:"account_number"`
	Sequence      int    `json:"sequence"`
}

type AccountResponseResult struct {
	Value *Account `json:"value"`
}

type AccountResponse struct {
	Result *AccountResponseResult `json:"result"`
}

func (ctx *Client) Account() (*Account, error) {
	res := &AccountResponse{}

	body, err := ctx.Query("/auth/accounts/" + ctx.Address)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res.Result.Value, nil
}

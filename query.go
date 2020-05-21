package bluzelle

import (
	"encoding/json"
	"fmt"
	"strconv"
)

//

type ReadResponseResult struct {
	Value string `json:"value"`
}

type ReadResponse struct {
	Result *ReadResponseResult `json:"result"`
}

func (ctx *Client) Read(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return "", err
	}

	body, err := ctx.APIQuery("/crud/read/" + ctx.options.UUID + "/" + encodeSafe(key))
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
	if key == "" {
		return "", fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return "", err
	}

	body, err := ctx.APIQuery("/crud/pread/" + ctx.options.UUID + "/" + encodeSafe(key))
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

//

type HasResponseResult struct {
	Has bool `json:"has"`
}

type HasResponse struct {
	Result *HasResponseResult `json:"result"`
}

func (ctx *Client) Has(key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return false, err
	}

	body, err := ctx.APIQuery("/crud/has/" + ctx.options.UUID + "/" + encodeSafe(key))
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

//

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

//

type KeysResponseResult struct {
	Keys []string `json:"keys"`
}

type KeysResponse struct {
	Result *KeysResponseResult `json:"result"`
}

func (ctx *Client) Keys() ([]string, error) {
	body, err := ctx.APIQuery("/crud/keys/" + ctx.options.UUID)
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

//

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

//

type GetLeaseResponseResult struct {
	Lease string `json:"lease"`
}

type GetLeaseResponse struct {
	Result *GetLeaseResponseResult `json:"result"`
}

func (ctx *Client) GetLease(key string) (int64, error) {
	if key == "" {
		return 0, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return 0, err
	}

	body, err := ctx.APIQuery("/crud/getlease/" + ctx.options.UUID + "/" + encodeSafe(key))
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

//

type GetNShortestLeasesResponseResultKeyLease struct {
	Key   string `json:"key"`
	Lease int64  `json:"lease"`
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
	body, err := ctx.APIQuery(fmt.Sprintf("/crud/getnshortestleases/%s/%d", ctx.options.UUID, n))
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

//

type Account struct {
	AccountNumber int     `json:"account_number"`
	Address       string  `json:"address"`
	Coins         []*Coin `json:"coins"`
	PublicKey     string  `json:"public_key"`
	Sequence      int     `json:"sequence"`
}

type Coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type AccountResponseResult struct {
	Value *Account `json:"value"`
}

type AccountResponse struct {
	Result *AccountResponseResult `json:"result"`
}

func (ctx *Client) Account() (*Account, error) {
	res := &AccountResponse{}

	body, err := ctx.APIQuery("/auth/accounts/" + ctx.Address)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res.Result.Value, nil
}

//

type VersionResponseApplicationVersion struct {
	Version string `json:"version"`
}

type VersionResponse struct {
	ApplicationVersion *VersionResponseApplicationVersion `json:"application_version"`
}

func (ctx *Client) Version() (string, error) {
	body, err := ctx.APIQuery("/node_info")
	if err != nil {
		return "", err
	}

	res := &VersionResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return "", err
	}
	return res.ApplicationVersion.Version, nil
}

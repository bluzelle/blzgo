package bluzelle

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const INVALID_LEASE_TIME = "Invalid lease time"

//

func (ctx *Client) Create(key string, value string, gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
	}
	if lease < 0 {
		return fmt.Errorf("%s", INVALID_LEASE_TIME)
	}

	transaction := &Transaction{
		Key:                key,
		Value:              value,
		Lease:              lease,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/create",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

//

func (ctx *Client) Update(key string, value string, gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
	}
	if lease < 0 {
		return fmt.Errorf("%s", INVALID_LEASE_TIME)
	}

	transaction := &Transaction{
		Key:                key,
		Value:              value,
		Lease:              lease,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/update",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

//

func (ctx *Client) Delete(key string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		ApiRequestMethod:   "DELETE",
		ApiRequestEndpoint: "/crud/delete",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

//

func (ctx *Client) Rename(key string, newKey string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		NewKey:             newKey,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/rename",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

//

func (ctx *Client) DeleteAll(gasInfo *GasInfo) error {
	transaction := &Transaction{
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/deleteall",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

//

func (ctx *Client) MultiUpdate(keyValues []*KeyValue, gasInfo *GasInfo) error {
	transaction := &Transaction{
		KeyValues:          keyValues,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/multiupdate",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

//

func (ctx *Client) RenewLease(key string, gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	if leaseInfo == nil {
		return fmt.Errorf("lease is required")
	}
	if lease < 0 {
		return fmt.Errorf("%s", INVALID_LEASE_TIME)
	}

	transaction := &Transaction{
		Key:                key,
		Lease:              leaseInfo.ToBlocks(),
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/renewlease",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

//

func (ctx *Client) RenewLeaseAll(gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	if leaseInfo == nil {
		return fmt.Errorf("lease is required")
	}
	if lease < 0 {
		return fmt.Errorf("%s", INVALID_LEASE_TIME)
	}

	transaction := &Transaction{
		Lease:              leaseInfo.ToBlocks(),
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/renewleaseall",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Client) RenewAllLeases(gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	return ctx.RenewLeaseAll(gasInfo, leaseInfo)
}

//

func (ctx *Client) TxRead(key string, gasInfo *GasInfo) (string, error) {
	transaction := &Transaction{
		Key:                key,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/read",
		GasInfo:            gasInfo,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return "", err
	}

	res := &ReadResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return "", err
	}
	return res.Value, nil
}

//

func (ctx *Client) TxHas(key string, gasInfo *GasInfo) (bool, error) {
	transaction := &Transaction{
		Key:                key,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/has",
		GasInfo:            gasInfo,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return false, err
	}

	res := &HasResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return false, err
	}
	return res.Has, nil
}

//
func (ctx *Client) TxCount(gasInfo *GasInfo) (int, error) {
	transaction := &Transaction{
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/count",
		GasInfo:            gasInfo,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return 0, err
	}

	res := &CountResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return 0, err
	}

	if count, err := strconv.Atoi(res.Count); err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

//

func (ctx *Client) TxKeys(gasInfo *GasInfo) ([]string, error) {
	transaction := &Transaction{
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/keys",
		GasInfo:            gasInfo,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return nil, err
	}

	res := &KeysResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

//

func (ctx *Client) TxKeyValues(gasInfo *GasInfo) ([]*KeyValue, error) {
	transaction := &Transaction{
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/keyvalues",
		GasInfo:            gasInfo,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return nil, err
	}

	res := &KeyValuesResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res.KeyValues, nil
}

//

func (ctx *Client) TxGetLease(key string, gasInfo *GasInfo) (int64, error) {
	transaction := &Transaction{
		Key:                key,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/getlease",
		GasInfo:            gasInfo,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return 0, err
	}

	res := &GetLeaseResponseResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return 0, err
	}
	lease, err := strconv.Atoi(res.Lease)
	return int64(lease * BLOCK_TIME_IN_SECONDS), err
}

//

func (ctx *Client) TxGetNShortestLeases(n uint64, gasInfo *GasInfo) ([]*GetNShortestLeasesResponseResultKeyLease, error) {
	transaction := &Transaction{
		N:                  n,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/getnshortestlease",
		GasInfo:            gasInfo,
	}

	body, err := ctx.SendTransaction(transaction)
	if err != nil {
		return nil, err
	}

	res := &GetNShortestLeasesResponseResult{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	kl, err := res.GetHumanizedKeyLeases()
	return kl, err
}

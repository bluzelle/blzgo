package bluzelle

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const INVALID_LEASE_TIME = "Invalid lease time"

func (ctx *Client) Create(key string, value string, gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	if msg, err := ctx.create(key, value, leaseInfo); err != nil {
		return err
	} else {
		transaction.Messages = append(transaction.Messages, msg)
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Client) Update(key string, value string, gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	if key == "" {
		return fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return err
	}
	if value == "" {
		return fmt.Errorf(VALUE_IS_REQUIRED)
	}

	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
	}
	if lease < 0 {
		return fmt.Errorf("%s", INVALID_LEASE_TIME)
	}

	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type: "crud/update",
		Value: &TransactionMessageValue{
			Key:   key,
			Value: value,
			Lease: strconv.FormatInt(lease, 10),
		},
	})

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Client) Delete(key string, gasInfo *GasInfo) error {
	if key == "" {
		return fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return err
	}

	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type: "crud/delete",
		Value: &TransactionMessageValue{
			Key: key,
		},
	})

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Client) Rename(key string, newKey string, gasInfo *GasInfo) error {
	if key == "" {
		return fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return err
	}
	if newKey == "" {
		return fmt.Errorf(NEW_KEY_IS_REQUIRED)
	}
	if err := validateKey(newKey); err != nil {
		return err
	}

	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type: "crud/rename",
		Value: &TransactionMessageValue{
			Key:    key,
			NewKey: newKey,
		},
	})

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Client) DeleteAll(gasInfo *GasInfo) error {
	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type:  "crud/deleteall",
		Value: &TransactionMessageValue{},
	})

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Client) MultiUpdate(keyValues []*KeyValue, gasInfo *GasInfo) error {
	if len(keyValues) == 0 {
		return fmt.Errorf(KEY_VALUES_ARE_REQUIRED)
	}
	for _, kv := range keyValues {
		key := kv.Key
		if key == "" {
			return fmt.Errorf(KEY_IS_REQUIRED)
		}
		if err := validateKey(key); err != nil {
			return err
		}
	}

	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type: "crud/multiupdate",
		Value: &TransactionMessageValue{
			KeyValues: keyValues,
		},
	})

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Client) RenewLease(key string, gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	if key == "" {
		return fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return err
	}

	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
	}
	if lease < 0 {
		return fmt.Errorf("%s", INVALID_LEASE_TIME)
	}

	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type: "crud/renewlease",
		Value: &TransactionMessageValue{
			Key:   key,
			Lease: strconv.FormatInt(lease, 10),
		},
	})

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Client) RenewLeaseAll(gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
	}
	if lease < 0 {
		return fmt.Errorf("%s", INVALID_LEASE_TIME)
	}

	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type: "crud/renewleaseall",
		Value: &TransactionMessageValue{
			Lease: strconv.FormatInt(lease, 10),
		},
	})

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Client) RenewAllLeases(gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	return ctx.RenewLeaseAll(gasInfo, leaseInfo)
}

func (ctx *Client) TxRead(key string, gasInfo *GasInfo) (string, error) {
	if key == "" {
		return "", fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return "", err
	}

	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type: "crud/read",
		Value: &TransactionMessageValue{
			Key: key,
		},
	})

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

func (ctx *Client) TxHas(key string, gasInfo *GasInfo) (bool, error) {
	if key == "" {
		return false, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return false, err
	}

	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type: "crud/has",
		Value: &TransactionMessageValue{
			Key: key,
		},
	})

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

func (ctx *Client) TxCount(gasInfo *GasInfo) (int, error) {
	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type:  "crud/count",
		Value: &TransactionMessageValue{},
	})

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

func (ctx *Client) TxKeys(gasInfo *GasInfo) ([]string, error) {
	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type:  "crud/keys",
		Value: &TransactionMessageValue{},
	})

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

func (ctx *Client) TxKeyValues(gasInfo *GasInfo) ([]*KeyValue, error) {
	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type:  "crud/ukeyvalues",
		Value: &TransactionMessageValue{},
	})

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

func (ctx *Client) TxGetLease(key string, gasInfo *GasInfo) (int64, error) {
	if key == "" {
		return 0, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return 0, err
	}

	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type: "crud/getlease",
		Value: &TransactionMessageValue{
			Key: key,
		},
	})

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

func (ctx *Client) TxGetNShortestLeases(n uint64, gasInfo *GasInfo) ([]*GetNShortestLeasesResponseResultKeyLease, error) {
	transaction := &Transaction{
		GasInfo: gasInfo,
	}

	transaction.Messages = append(transaction.Messages, &TransactionMessage{
		Type: "crud/getnshortestleases",
		Value: &TransactionMessageValue{
			N: strconv.FormatUint(n, 10),
		},
	})

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

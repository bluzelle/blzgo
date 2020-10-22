package bluzelle

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (ctx *Client) create(key string, value string, leaseInfo *LeaseInfo) (*TransactionMessage, error) {
	if key == "" {
		return nil, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}
	if value == "" {
		return nil, fmt.Errorf(VALUE_IS_REQUIRED)
	}
	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
	}
	if lease < 0 {
		return nil, fmt.Errorf(INVALID_LEASE_TIME)
	}

	return &TransactionMessage{
		Type: "crud/create",
		Value: &TransactionMessageValue{
			Key:   key,
			Value: value,
			Lease: strconv.FormatInt(lease, 10),
		},
	}, nil
}

func (ctx *Client) update(key string, value string, leaseInfo *LeaseInfo) (*TransactionMessage, error) {
	if key == "" {
		return nil, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}
	if value == "" {
		return nil, fmt.Errorf(VALUE_IS_REQUIRED)
	}

	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
	}
	if lease < 0 {
		return nil, fmt.Errorf("%s", INVALID_LEASE_TIME)
	}

	return &TransactionMessage{
		Type: "crud/update",
		Value: &TransactionMessageValue{
			Key:   key,
			Value: value,
			Lease: strconv.FormatInt(lease, 10),
		},
	}, nil
}

func (ctx *Client) delete(key string) (*TransactionMessage, error) {
	if key == "" {
		return nil, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}

	return &TransactionMessage{
		Type: "crud/delete",
		Value: &TransactionMessageValue{
			Key: key,
		},
	}, nil
}

func (ctx *Client) rename(key string, newKey string) (*TransactionMessage, error) {
	if key == "" {
		return nil, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}
	if newKey == "" {
		return nil, fmt.Errorf(NEW_KEY_IS_REQUIRED)
	}
	if err := validateKey(newKey); err != nil {
		return nil, err
	}

	return &TransactionMessage{
		Type: "crud/rename",
		Value: &TransactionMessageValue{
			Key:    key,
			NewKey: newKey,
		},
	}, nil
}

func (ctx *Client) deleteAll() (*TransactionMessage, error) {
	return &TransactionMessage{
		Type:  "crud/deleteall",
		Value: &TransactionMessageValue{},
	}, nil
}

func (ctx *Client) multiUpdate(keyValues []*KeyValue) (*TransactionMessage, error) {
	if len(keyValues) == 0 {
		return nil, fmt.Errorf(KEY_VALUES_ARE_REQUIRED)
	}
	for _, kv := range keyValues {
		key := kv.Key
		if key == "" {
			return nil, fmt.Errorf(KEY_IS_REQUIRED)
		}
		if err := validateKey(key); err != nil {
			return nil, err
		}
	}

	return &TransactionMessage{
		Type: "crud/multiupdate",
		Value: &TransactionMessageValue{
			KeyValues: keyValues,
		},
	}, nil
}

func (ctx *Client) renewLease(key string, leaseInfo *LeaseInfo) (*TransactionMessage, error) {
	if key == "" {
		return nil, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}

	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
	}
	if lease < 0 {
		return nil, fmt.Errorf("%s", INVALID_LEASE_TIME)
	}

	return &TransactionMessage{
		Type: "crud/renewlease",
		Value: &TransactionMessageValue{
			Key:   key,
			Lease: strconv.FormatInt(lease, 10),
		},
	}, nil
}

func (ctx *Client) renewLeaseAll(leaseInfo *LeaseInfo) (*TransactionMessage, error) {
	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
	}
	if lease < 0 {
		return nil, fmt.Errorf("%s", INVALID_LEASE_TIME)
	}

	return &TransactionMessage{
		Type: "crud/renewleaseall",
		Value: &TransactionMessageValue{
			Lease: strconv.FormatInt(lease, 10),
		},
	}, nil
}

func (ctx *Client) renewAllLeases(leaseInfo *LeaseInfo) (*TransactionMessage, error) {
	return ctx.renewLeaseAll(leaseInfo)
}

func (ctx *Client) txRead(key string) (*TransactionMessage, error) {
	if key == "" {
		return nil, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}

	return &TransactionMessage{
		Type: "crud/read",
		Value: &TransactionMessageValue{
			Key: key,
		},
	}, nil
}

func (ctx *Client) parseTxReadResponse(body []byte) (string, error) {
	res := &ReadResponseResult{}
	if err := json.Unmarshal(body, res); err != nil {
		return "", err
	}
	return res.Value, nil
}

func (ctx *Client) txHas(key string) (*TransactionMessage, error) {
	if key == "" {
		return nil, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}

	return &TransactionMessage{
		Type: "crud/has",
		Value: &TransactionMessageValue{
			Key: key,
		},
	}, nil
}

func (ctx *Client) txCount() (*TransactionMessage, error) {
	return &TransactionMessage{
		Type:  "crud/count",
		Value: &TransactionMessageValue{},
	}, nil
}

func (ctx *Client) txKeys() (*TransactionMessage, error) {
	return &TransactionMessage{
		Type:  "crud/keys",
		Value: &TransactionMessageValue{},
	}, nil
}

func (ctx *Client) txKeyValues() (*TransactionMessage, error) {
	return &TransactionMessage{
		Type:  "crud/ukeyvalues",
		Value: &TransactionMessageValue{},
	}, nil
}

func (ctx *Client) txGetLease(key string) (*TransactionMessage, error) {
	if key == "" {
		return nil, fmt.Errorf(KEY_IS_REQUIRED)
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}

	return &TransactionMessage{
		Type: "crud/getlease",
		Value: &TransactionMessageValue{
			Key: key,
		},
	}, nil
}

func (ctx *Client) txGetNShortestLeases(n uint64) (*TransactionMessage, error) {
	return &TransactionMessage{
		Type: "crud/getnshortestleases",
		Value: &TransactionMessageValue{
			N: strconv.FormatUint(n, 10),
		},
	}, nil
}

package bluzelle

import (
	"fmt"
	"strconv"
)

func (ctx *Client) Transaction() *BatchTransaction {
	txn := &BatchTransaction{
		client: ctx,
	}
	return txn
}

type BatchTransaction struct {
	messages []*TransactionMessage
	client   *Client
	gasInfo  *GasInfo
}

func (ctx *BatchTransaction) Create(key string, value string, gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
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
		return fmt.Errorf(INVALID_LEASE_TIME)
	}

	ctx.messages = append(ctx.messages, &TransactionMessage{
		Type: "crud/create",
		Value: &TransactionMessageValue{
			Key:   key,
			Value: value,
			Lease: strconv.FormatInt(lease, 10),
		},
	})

	if ctx.gasInfo == nil {
		ctx.gasInfo = gasInfo
	} else {
		if gasInfo.MaxFee > ctx.gasInfo.MaxFee {
			ctx.gasInfo.MaxFee = gasInfo.MaxFee
		}
		if gasInfo.MaxGas > ctx.gasInfo.MaxGas {
			ctx.gasInfo.MaxGas = gasInfo.MaxGas
		}
		if gasInfo.GasPrice > ctx.gasInfo.GasPrice {
			ctx.gasInfo.GasPrice = gasInfo.GasPrice
		}
	}

	return nil
}

func (ctx *BatchTransaction) Execute() error {
	txn := &Transaction{
		GasInfo: ctx.gasInfo,
	}

	txn.Messages = append(txn.Messages, ctx.messages...)

	_, err := ctx.client.SendTransaction(txn)
	if err != nil {
		return err
	}
	return nil
}

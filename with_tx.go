package bluzelle

import "fmt"

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
	result   []byte
}

func (ctx *BatchTransaction) Create(key string, value string, gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	if msg, err := ctx.client.create(key, value, leaseInfo); err != nil {
		return err
	} else {
		ctx.messages = append(ctx.messages, msg)
	}

	ctx.updateGas(gasInfo)

	return nil
}

func (ctx *BatchTransaction) TxRead(key string, gasInfo *GasInfo) error {
	if msg, err := ctx.client.txRead(key); err != nil {
		return err
	} else {
		ctx.messages = append(ctx.messages, msg)
	}

	ctx.updateGas(gasInfo)

	return nil
}

func (ctx *BatchTransaction) updateGas(gasInfo *GasInfo) {
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
}

func (ctx *BatchTransaction) Execute() ([]byte, error) {
	if ctx.result {
		return nil, fmt.Errorf("transaction was executed")
	}
	// if ctx.executing {
	// 	return nil, fmt.Errorf("transaction was executed")
	// }

	txn := &Transaction{
		GasInfo: ctx.gasInfo,
	}

	txn.Messages = append(txn.Messages, ctx.messages...)

	if result, err := ctx.client.SendTransaction(txn); err != nil {
		return err
	} else {
		ctx.result = result
		return nil
	}
}

func (ctx *BatchTransactionResult) GetTxReadResult(index int) (string, error) {
	// return ctx.client.txRead(ctx.result)
	return "", nil
}

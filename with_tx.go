package bluzelle

// import "fmt"

// func (ctx *Client) Transaction(uuid string) *MessageTransaction {
// 	txn := &MessageTransaction{}

// 	return txn
// }

// type MessageTransaction struct {
// }

// func (ctx *Client) MessageTransaction(key string, value string, gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
// 	if key == "" {
// 		return fmt.Errorf(KEY_IS_REQUIRED)
// 	}
// 	if err := validateKey(key); err != nil {
// 		return err
// 	}
// 	if value == "" {
// 		return fmt.Errorf(VALUE_IS_REQUIRED)
// 	}
// 	var lease int64
// 	if leaseInfo != nil {
// 		lease = leaseInfo.ToBlocks()
// 	}
// 	if lease < 0 {
// 		return fmt.Errorf(INVALID_LEASE_TIME)
// 	}

// 	transaction := &Transaction{
// 		Key:                key,
// 		Value:              value,
// 		Lease:              lease,
// 		ApiRequestMethod:   "POST",
// 		ApiRequestEndpoint: "/crud/create",
// 		GasInfo:            gasInfo,
// 	}

// 	_, err := ctx.SendTransaction(transaction)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (ctx *MessageTransaction) Execute() error {
// 	if key == "" {
// 		return fmt.Errorf(KEY_IS_REQUIRED)
// 	}
// 	if err := validateKey(key); err != nil {
// 		return err
// 	}
// 	if value == "" {
// 		return fmt.Errorf(VALUE_IS_REQUIRED)
// 	}
// 	var lease int64
// 	if leaseInfo != nil {
// 		lease = leaseInfo.ToBlocks()
// 	}
// 	if lease < 0 {
// 		return fmt.Errorf(INVALID_LEASE_TIME)
// 	}

// 	transaction := &Transaction{
// 		Key:                key,
// 		Value:              value,
// 		Lease:              lease,
// 		ApiRequestMethod:   "POST",
// 		ApiRequestEndpoint: "/crud/create",
// 		GasInfo:            gasInfo,
// 	}

// 	_, err := ctx.SendTransaction(transaction)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

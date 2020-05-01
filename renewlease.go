package bluzelle

import (
	"fmt"
)

func (ctx *Client) RenewLease(key string, leaseInfo *LeaseInfo, gasInfo *GasInfo) error {
	if leaseInfo == nil {
		return fmt.Errorf("lease is required")
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

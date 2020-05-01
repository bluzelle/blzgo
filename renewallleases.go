package bluzelle

import (
	"fmt"
)

func (ctx *Client) RenewAllLeases(leaseInfo *LeaseInfo, gasInfo *GasInfo) error {
	if leaseInfo == nil {
		return fmt.Errorf("lease is required")
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

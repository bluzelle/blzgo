package bluzelle

func (ctx *Client) Update(key string, value string, leaseInfo *LeaseInfo, gasInfo *GasInfo) error {
	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
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

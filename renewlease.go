package bluzelle

func (ctx *Client) RenewLease(key string, lease int64, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		Lease:              lease,
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

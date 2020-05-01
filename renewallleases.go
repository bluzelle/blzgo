package bluzelle

func (ctx *Client) RenewAllLeases(lease int64, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Lease:              lease,
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

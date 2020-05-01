package bluzelle

func (ctx *Client) Create(key string, value string, lease int64, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		Value:              value,
		Lease:              lease,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/create",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

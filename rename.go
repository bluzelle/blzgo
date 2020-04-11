package bluzelle

func (ctx *Client) Rename(key string, newKey string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		NewKey:             newKey,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/rename",
		GasInfo:            gasInfo,
		Client:             ctx,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

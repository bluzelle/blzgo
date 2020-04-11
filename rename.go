package bluzelle

func (ctx *Client) Rename(key string, newKey string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		NewKey:             newKey,
		Address:            ctx.Options.Address,
		UUID:               ctx.Options.UUID,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/rename",
		GasInfo:            gasInfo,
		ChainId:            ctx.Options.ChainId,
		Client:             ctx,
	}

	err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

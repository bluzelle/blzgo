package bluzelle

func (ctx *Client) Create(key string, value string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		Value:              value,
		Address:            ctx.Options.Address,
		UUID:               ctx.Options.UUID,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/create",
		GasInfo:            gasInfo,
		ChainId:            ctx.Options.ChainId,
		Client:             ctx,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

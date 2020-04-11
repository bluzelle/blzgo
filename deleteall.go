package bluzelle

func (ctx *Client) DeleteAll(uuid string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Address:            ctx.Options.Address,
		UUID:               uuid,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/deleteall",
		GasInfo:            gasInfo,
		ChainId:            ctx.Options.ChainId,
		Client:             ctx,
	}

	_, err := ctx.SendTransaction(transaction)
	return err
}

package bluzelle

func (ctx *Client) MultiUpdate(keyValues []*KeyValue, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Address:            ctx.Options.Address,
		UUID:               ctx.Options.UUID,
		KeyValues:          keyValues,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/multiupdate",
		GasInfo:            gasInfo,
		ChainId:            ctx.Options.ChainId,
		Client:             ctx,
	}

	_, err := ctx.SendTransaction(transaction)
	return err
}

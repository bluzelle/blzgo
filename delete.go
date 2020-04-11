package bluzelle

func (ctx *Client) Delete(key string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		Address:            ctx.Options.Address,
		UUID:               ctx.Options.UUID,
		ApiRequestMethod:   "DELETE",
		ApiRequestEndpoint: "/crud/delete",
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

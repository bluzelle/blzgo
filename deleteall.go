package bluzelle

func (ctx *Client) DeleteAll(gasInfo *GasInfo) error {
	transaction := &Transaction{
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/deleteall",
		GasInfo:            gasInfo,
		Client:             ctx,
	}

	_, err := ctx.SendTransaction(transaction)
	return err
}

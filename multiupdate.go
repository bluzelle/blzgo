package bluzelle

func (ctx *Client) MultiUpdate(keyValues []*KeyValue, gasInfo *GasInfo) error {
	transaction := &Transaction{
		KeyValues:          keyValues,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/multiupdate",
		GasInfo:            gasInfo,
		Client:             ctx,
	}

	_, err := ctx.SendTransaction(transaction)
	return err
}

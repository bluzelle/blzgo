package bluzelle

func (ctx *Client) MultiUpdate(keyValues []*KeyValue, gasInfo *GasInfo) error {
	transaction := &Transaction{
		KeyValues:          keyValues,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/multiupdate",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

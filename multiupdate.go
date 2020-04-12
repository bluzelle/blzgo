package bluzelle

func (ctx *Client) MultiUpdate(keyValues []*KeyValue) error {
	transaction := &Transaction{
		KeyValues:          keyValues,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/multiupdate",
		Client:             ctx,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

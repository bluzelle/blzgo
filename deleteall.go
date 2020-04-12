package bluzelle

func (ctx *Client) DeleteAll() error {
	transaction := &Transaction{
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/deleteall",
		Client:             ctx,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

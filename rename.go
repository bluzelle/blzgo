package bluzelle

func (ctx *Client) Rename(key string, newKey string) error {
	transaction := &Transaction{
		Key:                key,
		NewKey:             newKey,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/rename",
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

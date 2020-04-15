package bluzelle

func (ctx *Client) Update(key string, value string) error {
	transaction := &Transaction{
		Key:                key,
		Value:              value,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/update",
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

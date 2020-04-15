package bluzelle

func (ctx *Client) Create(key string, value string) error {
	transaction := &Transaction{
		Key:                key,
		Value:              value,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/create",
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

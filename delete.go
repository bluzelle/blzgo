package bluzelle

func (ctx *Client) Delete(key string) error {
	transaction := &Transaction{
		Key:                key,
		ApiRequestMethod:   "DELETE",
		ApiRequestEndpoint: "/crud/delete",
		Client:             ctx,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

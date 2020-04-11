package bluzelle

func (ctx *Client) Delete(key string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		ApiRequestMethod:   "DELETE",
		ApiRequestEndpoint: "/crud/delete",
		GasInfo:            gasInfo,
		Client:             ctx,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

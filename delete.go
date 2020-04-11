package bluzelle

type DeleteResponseResult struct {
	Key  string `json:"key"`
	UUID string `json:"uuid"`
}

type DeleteResponse struct {
	Result *DeleteResponseResult `json:"result"`
}

func (ctx *Client) Delete(key string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		Address:            ctx.Options.Address,
		UUID:               ctx.Options.UUID,
		ApiRequestMethod:   "DELETE",
		ApiRequestEndpoint: "/crud/delete",
		GasInfo:            gasInfo,
		ChainId:            ctx.Options.ChainId,
		Client:             ctx,
	}

	err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

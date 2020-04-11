package bluzelle

type UpdateResponseResult struct {
	Value string `json:"value"`
	Key   string `json:"key"`
	UUID  string `json:"uuid"`
}

type UpdateResponse struct {
	Result *UpdateResponseResult `json:"result"`
}

func (ctx *Client) Update(key string, value string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		Value:              value,
		Address:            ctx.Options.Address,
		UUID:               ctx.Options.UUID,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/update",
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

package bluzelle

type RenameResponseResult struct {
	Value string `json:"value"`
	Key   string `json:"key"`
	UUID  string `json:"uuid"`
}

type RenameResponse struct {
	Result *RenameResponseResult `json:"result"`
}

func (ctx *Client) Rename(key string, newKey string, gasInfo *GasInfo) error {
	transaction := &Transaction{
		Key:                key,
		NewKey:             newKey,
		Address:            ctx.Options.Address,
		UUID:               ctx.Options.UUID,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/rename",
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

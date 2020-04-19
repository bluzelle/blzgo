package bluzelle

func (ctx *Client) RenewAllLeases(lease int64) error {
	transaction := &Transaction{
		Lease:              lease,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/renewleaseall",
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

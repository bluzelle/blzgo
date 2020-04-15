package bluzelle

func (ctx *Client) RenewLeaseAll(lease int64) error {
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

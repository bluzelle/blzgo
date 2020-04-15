package bluzelle

func (ctx *Client) RenewLease(key string, lease int64) error {
	transaction := &Transaction{
		Key:                key,
		Lease:              lease,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/renewlease",
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

package bluzelle

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func (ctx *Client) APIQuery(endpoint string) ([]byte, error) {
	url := ctx.Options.Endpoint + endpoint

	ctx.Infof("get %s", url)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (ctx *Client) APIMutate(method string, endpoint string, payload []byte) ([]byte, error) {
	url := ctx.Options.Endpoint + endpoint

	ctx.Infof("post %s", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (ctx *Client) SendTransaction(transaction *Transaction) error {
	transaction.done = make(chan error, 1)
	ctx.transactions <- transaction
	err := <-transaction.done
	if err != nil {
		ctx.Errorf("transaction err(%s)", err)
	}
	return err
}

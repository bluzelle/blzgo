package bluzelle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func (ctx *Client) APIQuery(endpoint string) ([]byte, error) {
	url := ctx.options.Endpoint + endpoint

	ctx.Infof("get %s", url)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := parseResponse(res)
	return body, nil
}

func (ctx *Client) APIMutate(method string, endpoint string, payload []byte) ([]byte, error) {
	url := ctx.options.Endpoint + endpoint

	ctx.Infof("post %s", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := parseResponse(res)
	return body, nil
}

func (ctx *Client) SendTransaction(transaction *Transaction) ([]byte, error) {
	transaction.done = make(chan bool, 1)
	ctx.transactions <- transaction
	done := <-transaction.done
	if !done {
		ctx.Fatalf("txn did not complete") // todo: enqueue
	}
	if transaction.err != nil {
		ctx.Errorf("transaction err(%s)", transaction.err)
	}
	return transaction.result, transaction.err
}

func parseResponse(res *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	errRes := &ErrorResponse{}
	err = json.Unmarshal(body, errRes)
	if err != nil {
		return nil, err
	}

	if errRes.Error != "" {
		return nil, fmt.Errorf("%s", errRes.Error)
	}

	return body, nil
}

package bluzelle

import (
	"github.com/apex/log"
	"io/ioutil"
	"net/http"
)

func (ctx *Client) Query(endpoint string) ([]byte, error) {
	url := ctx.Endpoint + endpoint

	log.Infof("get %s", url)

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

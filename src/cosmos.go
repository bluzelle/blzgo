package bluzelle

import (
	"github.com/apex/log"
	"io/ioutil"
	// "net/url"
	"net/http"
)

func (ctx *Client) Query(endpoint string) ([]byte, error) {
	get_url := ctx.Endpoint + endpoint

	log.Infof("get %s", get_url)

	res, err := http.Get(get_url)
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

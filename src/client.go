package bluzelle

const DEFAULT_ENDPOINT string = "http://localhost:1317"

type ClientOptions struct {
	Address  string
	Mnemonic string
	Endpoint string
	UUID     string
	ChainId  string
	Debug    bool
}

type Client struct {
	Address  string
	Mnemonic string
	Endpoint string
	UUID     string
	ChainId  string
	Debug    bool
}

func NewClient(options *ClientOptions) (*Client, error) {
	client := &Client{
		Address:  options.Address,
		Mnemonic: options.Mnemonic,
		UUID:     options.UUID,
		Endpoint: options.Endpoint,
		ChainId:  options.ChainId,
	}

	if client.Endpoint == "" {
		client.Endpoint = DEFAULT_ENDPOINT
	}

	// generate private key from mnemonic

	// validate address against mnemonic

	// get account number and sequence

	return client, nil
}

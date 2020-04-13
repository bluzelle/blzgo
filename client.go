package bluzelle

import (
	"github.com/apex/log"
	"github.com/btcsuite/btcd/btcec"
)

const DEFAULT_ENDPOINT string = "http://localhost:1317"
const DEFAULT_CHAIN_ID string = "bluzelle"
const HD_PATH string = "m/44'/118'/0'/0/0"
const ADDRESS_PREFIX string = "bluzelle"

type Options struct {
	Address  string
	Mnemonic string
	Endpoint string
	UUID     string
	ChainId  string
	GasInfo  *GasInfo
	Debug    bool
}

type Client struct {
	options      *Options
	account      *Account
	logger       *log.Entry
	privateKey   *btcec.PrivateKey
	transactions chan *Transaction
}

func (ctx *Client) SendTransactions() {
	for transaction := range ctx.transactions {
		// ctx.Infof("processing transaction(%+v)", transaction)
		transaction.Send()
	}
}

func (root *Client) UUID(uuid string) *Client {
	options := &Options{
		Address: root.options.Address,
		// Mnemonic: root.options.Mnemonic,
		Endpoint: root.options.Endpoint,
		UUID:     uuid,
		ChainId:  root.options.ChainId,
		GasInfo:  root.options.GasInfo,
		Debug:    root.options.Debug,
	}

	ctx := &Client{
		options:      options,
		account:      root.account,
		privateKey:   root.privateKey,
		transactions: root.transactions,
	}

	ctx.setupLogger()

	return ctx
}

func (ctx *Client) setupLogger() {
	ctx.logger = log.WithFields(log.Fields{
		"uuid":    ctx.options.UUID,
		"address": ctx.options.Address,
	})
}

func NewClient(options *Options) (*Client, error) {
	if options.Endpoint == "" {
		options.Endpoint = DEFAULT_ENDPOINT
	}

	if options.ChainId == "" {
		options.ChainId = DEFAULT_CHAIN_ID
	}

	if options.UUID == "" {
		options.UUID = options.Address
	}

	if options.GasInfo == nil {
		options.GasInfo = &GasInfo{ // todo
			MaxFee: 4000001,
		}
	}

	ctx := &Client{
		options: options,
	}

	ctx.setupLogger()

	// Generate private key from mnemonic and compute address
	if err := ctx.setPrivateKey(); err != nil {
		return nil, err
	}

	// Validate the address against mnemonic
	if err := ctx.verifyAddress(); err != nil {
		return nil, err
	}

	// Fetch the address account info (`number` and `sequence` to be used later)
	if account, err := ctx.ReadAccount(); err != nil {
		return nil, err
	} else {
		ctx.account = account
	}

	// Send transactions
	ctx.transactions = make(chan *Transaction, 1) // serial
	go ctx.SendTransactions()

	return ctx, nil
}

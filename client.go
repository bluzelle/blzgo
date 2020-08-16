package bluzelle

import (
	"fmt"

	"github.com/apex/log"
	"github.com/btcsuite/btcd/btcec"
)

const DEFAULT_ENDPOINT string = "http://localhost:1317"
const DEFAULT_CHAIN_ID string = "bluzelle"
const HD_PATH string = "m/44'/118'/0'/0/0"
const ADDRESS_PREFIX string = "bluzelle"

type Options struct {
	Mnemonic string
	Endpoint string
	UUID     string
	ChainId  string
	Debug    bool
}

type Client struct {
	Address          string
	options          *Options
	account          *Account
	logger           *log.Entry
	privateKey       *btcec.PrivateKey
	broadcastRetries int
	transactions     chan *Transaction
}

func (root *Client) UUID(uuid string) *Client {
	options := &Options{
		Endpoint: root.options.Endpoint,
		UUID:     uuid,
		ChainId:  root.options.ChainId,
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

func (root *Client) Transaction() *Client {
	options := &Options{
		Endpoint: root.options.Endpoint,
		UUID:     root.options.UUID,
		ChainId:  root.options.ChainId,
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
		"uuid": ctx.options.UUID,
	})
}

// Fetch the address account info (`number` and `sequence` to be used later)
func (ctx *Client) setAccount() error {
	if account, err := ctx.Account(); err != nil {
		return err
	} else {
		ctx.account = account
		return nil
	}
}

func (ctx *Client) processTransactions() {
	for txn := range ctx.transactions {
		// ctx.Infof("processing transaction(%+v)", txn)
		ctx.ProcessTransaction(txn)
	}
}

func NewClient(options *Options) (*Client, error) {
	if options.Mnemonic == "" {
		return nil, fmt.Errorf("mnemonic is required")
	}
	if options.UUID == "" {
		return nil, fmt.Errorf("uuid is required")
	}
	if options.Endpoint == "" {
		options.Endpoint = DEFAULT_ENDPOINT
	}
	if options.ChainId == "" {
		options.ChainId = DEFAULT_CHAIN_ID
	}

	ctx := &Client{
		options: options,
	}

	ctx.setupLogger()

	// Generate private key from mnemonic
	if err := ctx.setPrivateKey(); err != nil {
		return nil, err
	}

	// Derive address from the mnemonic
	if err := ctx.setAddress(); err != nil {
		return nil, err
	}

	// Fetch the address account info (`number` and `sequence` to be used later)
	if err := ctx.setAccount(); err != nil {
		return nil, err
	}

	// Send transactions
	ctx.transactions = make(chan *Transaction, 1) // serial
	go ctx.processTransactions()

	return ctx, nil
}

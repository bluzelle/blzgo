package bluzelle

import (
	"crypto/sha256"
	"fmt"
	"github.com/apex/log"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160"
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
		// ctx.Infof("processing op for key(%+v)", transaction)
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
	ctx.logger = log.WithFields(log.Fields{})
}

func (ctx *Client) serveTransactions() {
	ctx.transactions = make(chan *Transaction, 1) // serial
	go ctx.SendTransactions()
}

func (ctx *Client) setPrivateKey() error {
	// generate private key from mnemonic
	if key, err := getECPrivateKey(ctx.options.Mnemonic, ctx.options.Address); err != nil {
		return err
	} else {
		ctx.privateKey = key
	}

	// validate address against mnemonic
	pubkey := ctx.privateKey.PubKey()
	x := secp256k1.CompressPubkey(pubkey.X, pubkey.Y)
	b := hashRipemd160(hashSha256(x))

	if z, err := bech32.ConvertBits(b, 8, 5, true); err != nil {
		return err
	} else {
		if a, err := bech32.Encode(ADDRESS_PREFIX, z); err != nil {
			return err
		} else if ctx.options.Address != a {
			return fmt.Errorf("Bad credentials - verify your address and mnemonic")
		}
	}

	// get account number and sequence
	if account, err := ctx.ReadAccount(); err != nil {
		return err
	} else {
		ctx.account = account
	}

	return nil
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

	if err := ctx.setPrivateKey(); err != nil {
		return nil, err
	}

	ctx.serveTransactions()

	return ctx, nil
}

func getECPrivateKey(mnemonic string, address string) (*btcec.PrivateKey, error) {
	seed := bip39.NewSeed(mnemonic, "")
	dpath, err := accounts.ParseDerivationPath(HD_PATH)
	if err != nil {
		return nil, err
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	key := masterKey
	for _, n := range dpath {
		key, err = key.Child(n)
		if err != nil {
			return nil, err
		}
	}
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func hashSha256(s []byte) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func hashRipemd160(s []byte) []byte {
	h := ripemd160.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

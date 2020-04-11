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
	Debug    bool
}

type Client struct {
	Options      *Options
	transactions chan *Transaction
	Account      *Account
	PrivateKey   *btcec.PrivateKey
	Logger       *log.Entry
}

func (ctx *Client) SendTransactions() {
	for transaction := range ctx.transactions {
		ctx.Infof("processing op for key(%s)", transaction.Key)
		transaction.Send()
	}
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

	ctx := &Client{
		Options: options,
	}

	ctx.Logger = log.WithFields(log.Fields{})

	// generate private key from mnemonic
	if key, err := getECPrivateKey(ctx.Options.Mnemonic, ctx.Options.Address); err != nil {
		return nil, err
	} else {
		ctx.PrivateKey = key
	}

	// validate address against mnemonic
	pubkey := ctx.PrivateKey.PubKey()
	x := secp256k1.CompressPubkey(pubkey.X, pubkey.Y)
	b := hashRipemd160(hashSha256(x))

	if z, err := bech32.ConvertBits(b, 8, 5, true); err != nil {
		return nil, err
	} else {
		if a, err := bech32.Encode(ADDRESS_PREFIX, z); err != nil {
			return nil, err
		} else if ctx.Options.Address != a {
			return nil, fmt.Errorf("Bad credentials - verify your address and mnemonic")
		}
	}

	// get account number and sequence
	if account, err := ctx.ReadAccount(); err != nil {
		return nil, err
	} else {
		ctx.Account = account
	}

	// serve transactions
	ctx.transactions = make(chan *Transaction, 1) // serial
	go ctx.SendTransactions()

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

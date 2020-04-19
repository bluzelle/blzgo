package bluzelle

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tyler-smith/go-bip39"
)

// Generate private key from mnemonic and compute address
func (ctx *Client) setPrivateKey() error {
	// Get root key
	seed := bip39.NewSeed(ctx.options.Mnemonic, "")
	masterExtendedKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return err
	}

	// Find chain key at HD_PATH
	extendedKey := masterExtendedKey
	if dpath, err := accounts.ParseDerivationPath(HD_PATH); err != nil {
		return err
	} else {
		for _, n := range dpath {
			extendedKey, err = extendedKey.Child(n)
			if err != nil {
				return err
			}
		}
	}

	// Get private key from extended key
	if privateKey, err := extendedKey.ECPrivKey(); err != nil {
		return err
	} else {
		ctx.privateKey = privateKey
	}
	return nil
}

// Validate the address against mnemonic
func (ctx *Client) verifyAddress() error {
	var tmPubKey tmsecp256k1.PubKeySecp256k1
	copy(tmPubKey[:], ctx.privateKey.PubKey().SerializeCompressed()[:tmsecp256k1.PubKeySecp256k1Size])
	ctx.Warnf("%x", ctx.privateKey.PubKey().SerializeCompressed())
	// Get bech32 address
	if conv, err := bech32.ConvertBits(tmPubKey.Address().Bytes(), 8, 5, true); err != nil {
		return err
	} else {
		if encoded, err := bech32.Encode(ADDRESS_PREFIX, conv); err != nil {
			return err
		} else if ctx.options.Address != encoded {
			return fmt.Errorf("bad credentials(verify your address and mnemonic)")
		}
	}

	return nil
}

package main

import (
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo/src"
	"os"
)

func main() {
	setupLogging()
	loadEnv()

	// create client
	options := &bluzelle.ClientOptions{
		Address:  os.Getenv("ADDRESS"),
		Mnemonic: os.Getenv("MNEMONIC"),
		UUID:     os.Getenv("UUID"),
		Endpoint: os.Getenv("ENDPOINT"),
		ChainId:  os.Getenv("CHAIN_ID"),
		Debug:    false,
	}
	ctx, err := bluzelle.NewClient(options)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// read account
	if account, err := ctx.Account(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("account info: %+v", account)
	}
}

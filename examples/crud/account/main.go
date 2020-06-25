package main

import (
	"github.com/apex/log"
	bluzelle "github.com/bluzelle/blzgo"
)

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Infof("getting account info...")

	if account, err := ctx.Account(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("account info: %+v", account)
	}
}

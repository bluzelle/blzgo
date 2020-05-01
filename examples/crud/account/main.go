package main

import (
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo"
)

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Infof("getting account info...")

	if account, err := ctx.ReadAccount(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("account info: %+v", account)
	}
}

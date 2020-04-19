package main

import (
	"github.com/apex/log"
	util "github.com/vbstreetz/blzgo"
)

func main() {
	util.SetupLogging()
	util.LoadEnv()

	ctx, err := util.NewTestClient()
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

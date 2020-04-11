package main

import (
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo/examples/util"
)

func main() {
	util.SetupLogging()
	util.LoadEnv()

	ctx, err := util.NewClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Infof("getting account info for address(%s)...", ctx.Options.Address)

	if account, err := ctx.ReadAccount(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("account info: %+v", account)
	}
}

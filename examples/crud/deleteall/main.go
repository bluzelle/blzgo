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

	log.Infof("deleting all keys...")

	if err := ctx.DeleteAll(nil); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("done")
	}
}

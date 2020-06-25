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

	log.Infof("deleting all keys...")

	if err := ctx.DeleteAll(bluzelle.TestGasInfo()); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("done")
	}
}

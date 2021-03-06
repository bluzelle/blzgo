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

	log.Infof("getting keyvalues...")

	if keyValues, err := ctx.KeyValues(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("values:")
		for _, keyValue := range keyValues {
			log.Infof("key(%s) value(%+v)", keyValue.Key, keyValue.Value)
		}
	}
}

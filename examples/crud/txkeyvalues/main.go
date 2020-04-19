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

	log.Infof("getting keyvalues...")

	if keyValues, err := ctx.TxKeyValues(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("values:")
		for _, keyValue := range keyValues {
			log.Infof("key(%s) value(%+v)", keyValue.Key, keyValue.Value)
		}
	}
}

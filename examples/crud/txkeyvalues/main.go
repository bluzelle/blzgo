package main

import (
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo/examples/util"
	"os"
)

func main() {
	util.SetupLogging()
	util.LoadEnv()

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("uuid is required")
	}

	ctx, err := util.NewClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	uuid := args[0]

	log.Infof("getting uuid(%s) keyvalues...", uuid)

	if keyValues, err := ctx.TxKeyValues(uuid, util.GasInfo()); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("values:")
		for _, keyValue := range keyValues {
			log.Infof("key(%s) value(%+v)", keyValue.Key, keyValue.Value)
		}
	}
}

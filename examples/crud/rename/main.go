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
	if len(args) < 2 {
		log.Fatalf("both key and newkey are required")
	}

	ctx, err := util.NewClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]
	newKey := args[1]

	log.Infof("renaming key(%s) to new key(%s)...", key, newKey)

	if err := ctx.Rename(key, newKey, util.GasInfo()); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("renamed key")
	}
}

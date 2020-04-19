package main

import (
	"github.com/apex/log"
	util "github.com/vbstreetz/blzgo"
	"os"
)

func main() {
	util.SetupLogging()
	util.LoadEnv()

	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("both key and newkey are required")
	}

	ctx, err := util.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]
	newKey := args[1]

	log.Infof("renaming key(%s) to new key(%s)...", key, newKey)

	if err := ctx.Rename(key, newKey); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("renamed key")
	}
}

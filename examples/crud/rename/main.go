package main

import (
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo"
	"os"
)

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("both key and newkey are required")
	}

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]
	newKey := args[1]

	log.Infof("renaming key(%s) to new key(%s)...", key, newKey)

	if err := ctx.Rename(key, newKey, nil); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("renamed key")
	}
}

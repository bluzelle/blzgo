package main

import (
	"os"
	"strconv"

	"github.com/apex/log"
	bluzelle "github.com/bluzelle/blzgo"
)

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("key is required")
	}

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]
	prove := false

	if len(args) > 1 {
		if b, err := strconv.ParseBool(args[1]); err == nil {
			prove = b
		}
	}

	log.Infof("getting val for key(%s) prove(%t)...", key, prove)

	if v, err := ctx.Read(key); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("val(%s)", v)
	}
}

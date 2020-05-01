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
	if len(args) == 0 {
		log.Fatalf("key is required")
	}

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]

	log.Infof("getting val for key(%s)...", key)

	if v, err := ctx.TxRead(key, nil); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("val(%s)", v)
	}
}

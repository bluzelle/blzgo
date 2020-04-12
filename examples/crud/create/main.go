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
		log.Fatalf("both key and value are required")
	}

	ctx, err := util.NewClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]
	value := args[1]

	log.Infof("creating key(%s), val(%s)...", key, value)

	if err := ctx.Create(key, value); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("created key")
	}
}

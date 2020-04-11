package main

import (
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo/examples/util"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("key is required")
	}

	key := args[0]

	ctx, err := util.NewClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Infof("getting val for key(%s)...", key)

	if v, err := ctx.Read(key, false); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("val for key(%s): %s", key, v)
	}
}

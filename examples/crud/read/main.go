package main

import (
	"github.com/apex/log"
	util "github.com/vbstreetz/blzgo"
	"os"
	"strconv"
)

func main() {
	util.SetupLogging()
	util.LoadEnv()

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("key is required")
	}

	ctx, err := util.NewTestClient()
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

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
	if len(args) == 0 {
		log.Fatalf("key is required")
	}

	ctx, err := util.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]

	log.Infof("getting lease for key(%s)...", key)

	if v, err := ctx.TxGetLease(key); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("lease(%d)", v)
	}
}

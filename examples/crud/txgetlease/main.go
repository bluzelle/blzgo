package main

import (
	"os"

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

	log.Infof("getting lease for key(%s)...", key)

	if v, err := ctx.TxGetLease(key, bluzelle.TestGasInfo()); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("lease(%ds)", v)
	}
}

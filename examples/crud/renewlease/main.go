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
	if len(args) < 2 {
		log.Fatalf("both key and lease are required")
	}

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]
	leaseInfo := &bluzelle.LeaseInfo{}
	if lease, err := strconv.Atoi(args[1]); err != nil {
		log.Fatalf("could not parse provided lease(%s)", args[1])
	} else {
		leaseInfo.Seconds = int64(lease)
	}

	log.Infof("renewing key(%s), lease(%ds)...", key, leaseInfo.ToBlocks())

	if err := ctx.RenewLease(key, bluzelle.GetTestGasInfo(), leaseInfo); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("renewed lease")
	}
}

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
	if len(args) < 1 {
		log.Fatalf("lease is required")
	}

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	leaseInfo := &bluzelle.LeaseInfo{}
	if lease, err := strconv.Atoi(args[0]); err != nil {
		log.Fatalf("could not parse provided lease(%s)", args[0])
	} else {
		leaseInfo.Seconds = int64(lease)
	}

	log.Infof("renewing lease(%ds)...", leaseInfo.ToBlocks())

	if err := ctx.RenewAllLeases(bluzelle.TestGasInfo(), leaseInfo); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("renewed leases")
	}
}

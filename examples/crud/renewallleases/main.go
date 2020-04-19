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
	if len(args) < 1 {
		log.Fatalf("lease is required")
	}

	ctx, err := util.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	lease := 0
	if l, err := strconv.Atoi(args[0]); err != nil {
		log.Fatalf("could not parse provided lease(%s)", args[0])
	} else {
		lease = l
	}

	log.Infof("renewing leases(%d)...", lease)

	if err := ctx.RenewAllLeases(int64(lease)); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("renewed leases")
	}
}

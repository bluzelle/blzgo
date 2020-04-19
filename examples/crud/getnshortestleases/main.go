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
		log.Fatalf("n is required")
	}

	ctx, err := util.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	n := uint64(0)
	if i, err := strconv.ParseUint(args[0], 10, 64); err != nil {
		log.Fatalf("could not parse provided n(%d)", err)
	} else {
		n = i
	}

	log.Infof("getting n shortest leases n(%d)...", n)

	if keyLeases, err := ctx.GetNShortestLeases(uint64(n)); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("leases:")
		for _, keyLease := range keyLeases {
			log.Infof("key(%s) lease(%+v)", keyLease.Key, keyLease.Lease)
		}
	}
}

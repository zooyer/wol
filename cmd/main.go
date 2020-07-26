package main

import (
	"flag"
	"os"

	"github.com/zooyer/wol"
)

var mac = flag.String("mac", "", "mac address")

func init() {
	flag.Parse()

	if *mac == "" && len(os.Args) > 1 {
		*mac = os.Args[1]
	}

	if *mac == "" {
		flag.Usage()
		os.Exit(2)
	}
}

// 54-E1-AD-10-44-EA
func main() {
	var err error
	if err = wol.WOL(*mac); err != nil {
		panic(err)
	}
}

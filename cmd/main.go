package main

import (
	"flag"
	"os"

	"github.com/zooyer/wol"
)

var (
	port = 9
	addr string
	macs []string

	mac = flag.String("mac", "", "mac address")
)

var args struct {
	Port   []int    `clop:"-c;--port" usage:""`
	MAC    []string `clop:"-m;--mac" usage:""`
	Remote string   `clop:"-r;--remote" usage:""`
}

func init() {
	//if err := clop.Bind(&args); err != nil {
	//	panic(err)
	//}

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

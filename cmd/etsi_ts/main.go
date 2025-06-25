package main

import (
	"flag"
	"fmt"

	"github.com/SUNET/g119612/pkg/etsi119612"
)

var Version = "1.0.0"

var urlVar = ""
var skiVar = ""

func config() {
	flag.StringVar(&urlVar, "url", "", "The URL of the trust status list")
	flag.StringVar(&skiVar, "ski", "", "The public subject key identifier to look for")
}

func main() {

	config()
	flag.Parse()

	tsl, err := etsi119612.FetchTSL(urlVar)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	pool := tsl.ToCertPool(etsi119612.PolicyAll)
	fmt.Printf("pool %+v", pool)
}

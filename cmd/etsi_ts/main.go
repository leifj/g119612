package main

import (
	"crypto/x509"
	"encoding/base64"
	"flag"
	"fmt"

	"github.com/SUNET/g119612/pkg/etsi119612"
)

var Version = "1.0.0"

var urlVar = ""
var x5cVar = ""

func config() {
	flag.StringVar(&urlVar, "url", "", "URL of a trust status list")
	flag.StringVar(&x5cVar, "x5c", "", "base64 encoded certificate (single line)")
}

func main() {

	config()
	flag.Parse()

	tsl, err := etsi119612.FetchTSL(urlVar)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	data, err := base64.StdEncoding.DecodeString(x5cVar)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	cert, err := x509.ParseCertificate(data)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	pool := tsl.ToCertPool(etsi119612.PolicyAll)
	_, err = cert.Verify(x509.VerifyOptions{Roots: pool})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Print("OK!\n")
}

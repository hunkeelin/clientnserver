package main

import (
	"encoding/pem"
	"flag"
	"github.com/hunkeelin/pki"
	"os"
)

var (
	gconfdir = flag.String("config", "client.conf", "location of the genkins.conf")
)

func main() {
	h, _ := os.Hostname()
	flag.Parse()
	c := readconfig(*gconfdir)
	keypath := c.keycertdir + "keys/" + h + ".key"
	csr, key := klinpki.GenCSRv2(2048)
	keyOut, err := os.OpenFile(keypath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	pem.Encode(keyOut, key)
	masteraddr := getHostnameFromCert(c.mastercrt)
	url := "https://" + masteraddr + ":" + c.masterport
	f, err := sendcsr(c.mastercrt, url, c.keycertdir, csr.Bytes)
	if err != nil {
		panic(err)
	}
	clientCRTFile, err := os.Create(h + ".crt")
	if err != nil {
		panic(err)
	}
	pem.Encode(clientCRTFile, &pem.Block{Type: "CERTIFICATE", Bytes: f})
	clientCRTFile.Close()
}

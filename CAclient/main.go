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
	flag.Parse()
	c := readconfig(*gconfdir)
	keypath := c.keycertdir + "keys/test2.klin-pro.com.key"
	csr, key := klinpki.GenCSRv2(2048)
	keyOut, err := os.OpenFile(keypath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	pem.Encode(keyOut, key)
	masteraddr := getHostnameFromCert(c.mastercrt)
	url := "https://" + masteraddr + ":" + c.masterport
	sendcsr(c.mastercrt, url, c.keycertdir, csr.Bytes)
}

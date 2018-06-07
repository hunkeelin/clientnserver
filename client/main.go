package main

import (
	"flag"
	"os"
)

var (
	gconfdir = flag.String("config", "client.conf", "location of the genkins.conf")
)

func main() {
	flag.Parse()
	c := readconfig(*gconfdir)
	h, _ := os.Hostname()
	certpath := c.keycertdir + "certs/" + h + ".crt"
	keypath := c.keycertdir + "keys/" + h + ".key"
	if !exist(certpath) || !exist(keypath) {
		genCrt(certpath, keypath)
	}
	masteraddr := getHostnameFromCert(c.mastercrt)
	url := "https://" + masteraddr + ":" + c.masterport
	sendcert(c.mastercrt, url, c.keycertdir)
}

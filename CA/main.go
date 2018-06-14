package main

import (
	"flag"
	"github.com/hunkeelin/pki"
)

var (
	gconfdir = flag.String("config", "CA.conf", "location of the genkins.conf")
	genCa    = flag.Bool("genCA", false, "true = generate ca")
)

func main() {
	flag.Parse()
	c := readconfig(*gconfdir)
	if *genCa {
		j := &klinpki.PkiConfig{
			EmailAddress: "support" + c.org + ".com",
			EcdsaCurve:   "",
			Certpath:     c.CApath,
			Keypath:      c.CAkeypath,
			MaxDays:      7200,
			RsaBits:      4096,
			Organization: c.org,
		}
		klinpki.GenCA(j)
		return
	}
	runServer(&c)
}

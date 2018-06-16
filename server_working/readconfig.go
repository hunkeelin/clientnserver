package main

import (
	"github.com/hunkeelin/klinenv"
	"log"
	"strconv"
)

func readconfig(p string) Config {
	var c Config

	config := klinenv.NewAppConfig(p)
	rconcur, err := config.Get("concur")
	if err != nil {
		log.Fatal("unable to retrieve the value of concur check config file")
	}
	concur, err := strconv.Atoi(rconcur)
	if err != nil {
		log.Fatal("can't convert string to int for concur")
	}
	c.concur = concur
	apikey, err := config.Get("apikey")
	checkerr(err)
	c.apikey = apikey

	bindaddr, err := config.Get("bindaddr")
	checkerr(err)
	c.bindaddr = bindaddr

	port, err := config.Get("port")
	checkerr(err)
	c.port = port

	certpath, err := config.Get("certpath")
	checkerr(err)
	c.certpath = certpath

	keypath, err := config.Get("keypath")
	checkerr(err)
	c.keypath = keypath

	pkidir, err := config.Get("pkidir")
	checkerr(err)
	if string(pkidir[len(pkidir)-1]) != "/" {
		pkidir += "/"
	}
	c.pkidir = pkidir
	return c
}

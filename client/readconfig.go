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

	mastercrt, err := config.Get("mastercrt")
	checkerr(err)
	c.mastercrt = mastercrt

	keycertdir, err := config.Get("keycertdir")
	checkerr(err)
	if len(keycertdir) == 0 {
		log.Fatal("please input keycertdir in config")
	} else {
		if string(keycertdir[len(keycertdir)-1]) != "/" {
			keycertdir += "/"
		}
	}
	c.keycertdir = keycertdir

	masterport, err := config.Get("masterport")
	checkerr(err)
	c.masterport = masterport
	return c
}

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/hunkeelin/pki"
	"log"
	"net/http"
	//	"os"
	"time"
)

var (
	gconfdir = flag.String("config", "client.conf", "location of the genkins.conf")
)

func main() {
	flag.Parse()
	c := readconfig(*gconfdir)
	csrpath := c.keycertdir + "csr/"
	keypath := c.keycertdir + "keys/test2.klin-pro.com.key"
	if !exist(keypath) {
		klinpki.GenCSR(2048, keypath, csrpath)
		masteraddr := getHostnameFromCert(c.mastercrt)
		url := "https://" + masteraddr + ":" + c.masterport
		sendcsr(c.mastercrt, url, c.keycertdir)
	} else {
		masteraddr := getHostnameFromCert(c.mastercrt)
		url := "https://" + masteraddr + ":" + c.masterport
		sendcsrv2(c.mastercrt, url, c.keycertdir)
		log.Fatal("done")
		newcon := new(Conn)
		newcon.monorun = make(chan struct{}, 1)
		tlsconfig := &tls.Config{
			PreferServerCipherSuites: true,
			// Only use curves which have assembly implementations
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			MinVersion: tls.VersionTLS12,
		}
		con := http.NewServeMux()
		con.HandleFunc("/", newcon.handleWebHook)
		s := &http.Server{
			Addr:         c.bindaddr + ":" + c.port,
			TLSConfig:    tlsconfig,
			Handler:      con,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		fmt.Println("listening to " + c.bindaddr + " " + c.port)
		certpath := "program/certs/test2.klin-pro.com.crt"
		err := s.ListenAndServeTLS(certpath, keypath)
		//err := s.ListenAndServe()
		if err != nil {
			log.Fatal("can't listen and serve check port and binding addr", err)
		}
	}
}

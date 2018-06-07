package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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
		masteraddr := getHostnameFromCert(c.mastercrt)
		url := "https://" + masteraddr + ":" + c.masterport
		sendcert(c.mastercrt, url, c.keycertdir)
	} else {
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
		err := s.ListenAndServeTLS(certpath, keypath)
		//err := s.ListenAndServe()
		if err != nil {
			log.Fatal("can't listen and serve check port and binding addr", err)
		}
	}
}

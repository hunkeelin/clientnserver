package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func runServer(c *Config) {
	newcon := new(Conn)
	// define config params
	sema := make(chan struct{}, 1)
	newcon.monorun = sema
	newcon.apikey = c.apikey
	newcon.concur = c.concur
	newcon.pkidir = c.pkidir

	if !Exist(c.certpath) || !Exist(c.keypath) {
		log.Fatal("please generate csr and sign it and put it in the correct directory")
	}
	certBytes, err := ioutil.ReadFile(c.CApath)
	if err != nil {
		log.Fatalln("Unable to read crt", err)
	}

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
		log.Fatalln("Unable to add certificate to certificate pool")
	}

	tlsconfig := &tls.Config{
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		ClientCAs: clientCertPool,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		//	ClientAuth: tls.RequireAndVerifyClientCert,
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
	err = s.ListenAndServeTLS(c.certpath, c.keypath)
	//err := s.ListenAndServe()
	if err != nil {
		log.Fatal("can't listen and serve check port and binding addr", err)
	}
}

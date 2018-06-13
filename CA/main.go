package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/hunkeelin/pki"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	gconfdir = flag.String("config", "server.conf", "location of the genkins.conf")
	sign     = flag.String("sign", "", "name of the csr you are signing")
	genCa    = flag.Bool("genCA", false, "true = generate ca")
	genCSR   = flag.Bool("genCSR", false, "true = generate CSR")
)

func main() {
	flag.Parse()
	if *sign != "" {
		// sign cert
		cerout, _ := klinpki.SignCSR("program/CA/ca.crt", "program/CA/ca.key", "program/pending/"+*sign+".csr", 20)
		clientCRTFile, err := os.Create("program/signed/" + *sign + ".crt")
		if err != nil {
			panic(err)
		}
		pem.Encode(clientCRTFile, &pem.Block{Type: "CERTIFICATE", Bytes: cerout})
		clientCRTFile.Close()
		os.Exit(1)
	}
	if *genCa == true {
		klinpki.GenCA("support@klin-pro.com", "", "program/CA/ca.crt", "program/CA/ca.key", 7200, 2048)
		os.Exit(1)
	}
	newcon := new(Conn)
	// define config params
	c := readconfig(*gconfdir)
	sema := make(chan struct{}, 1)
	newcon.monorun = sema
	newcon.apikey = c.apikey
	newcon.concur = c.concur
	newcon.pkidir = c.pkidir

	if *genCSR == true {
		klinpki.GenCSR(2048, c.keypath, "program/pending/")
		os.Exit(1)
	}

	// Generate cert
	if !Exist(c.certpath) || !Exist(c.keypath) {
		log.Fatal("please generate csr and sign it and put it in the correct directory")
	}
	certBytes, err := ioutil.ReadFile("program/CA/ca.crt")
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

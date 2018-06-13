package main

import (
	"crypto/tls"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/hunkeelin/pki"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	gconfdir = flag.String("config", "server.conf", "location of the genkins.conf")
)

func main() {
	flag.Parse()
	newcon := new(Conn)
	// define config params
	c := readconfig(*gconfdir)
	sema := make(chan struct{}, 1)
	newcon.monorun = sema
	newcon.apikey = c.apikey
	newcon.concur = c.concur
	newcon.pkidir = c.pkidir

	// Generate cert
	if !Exist(c.certpath) || !Exist(c.keypath) {
		klinpki.GenCSR(2048, c.keypath, "")
		klinpki.GenCA("support@klin-pro.com", "", "ca.crt", "ca.key", 7200, 2048)
		cerout, err := klinpki.SignCSR("ca.crt", "ca.key", "test1.klin-pro.com.csr", 50)
		if err != nil {
			panic(err)
		}
		clientCRTFile, err := os.Create("program/certs/test1.klin-pro.com.crt")
		if err != nil {
			panic(err)
		}
		pem.Encode(clientCRTFile, &pem.Block{Type: "CERTIFICATE", Bytes: cerout})
		clientCRTFile.Close()
	}
	// sign cert
	cerout, _ := klinpki.SignCSR(c.certpath, c.keypath, "/home/bgops/files/golesson/clientnserver/server/program/pending/test2.klin-pro.com.csr", 20)
	clientCRTFile, err := os.Create("test2.klin-pro.com" + ".crt")
	if err != nil {
		panic(err)
	}
	pem.Encode(clientCRTFile, &pem.Block{Type: "CERTIFICATE", Bytes: cerout})
	clientCRTFile.Close()
	//    certBytes, err := ioutil.ReadFile("client/cert.pem")
	//    if err != nil {
	//        log.Fatalln("Unable to read cert.pem", err)
	//    }
	//
	//    clientCertPool := x509.NewCertPool()
	//    if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
	//        log.Fatalln("Unable to add certificate to certificate pool")
	//    }

	tlsconfig := &tls.Config{
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		//ClientAuth: tls.RequireAndVerifyClientCert,
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

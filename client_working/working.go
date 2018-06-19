package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Creates a new file upload http request with optional extra params
func testreq(mastercert, host string) {
	cert, err := tls.LoadX509KeyPair("test2.klin-pro.com.crt", "test2.klin-pro.com.key")
	if err != nil {
		log.Fatalln("Unable to load cert", err)
	}

	// Load our CA certificate
	clientCACert, err := ioutil.ReadFile("master.pem")
	if err != nil {
		log.Fatal("Unable to open cert", err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            clientCertPool,
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{
		Timeout:   500 * time.Millisecond,
		Transport: tr,
	}
	payload := Records{{C: "asdf", D: true}}

	encodepayload, _ := json.Marshal(payload)
	ebody := bytes.NewReader(encodepayload)
	req, err := http.NewRequest("POST", "https://test3.klin-pro.com:2018", ebody)
	if err != nil {
		log.Println("Unable to speak to our server", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body), string(resp.Status))
}

type Records []struct {
	C string `json:"content"`
	D bool   `json:"disabled"`
}

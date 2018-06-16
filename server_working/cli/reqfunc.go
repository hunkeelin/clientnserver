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

type Payload struct {
	s string
}

func dowork(m string) (string, int) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	certs, err := ioutil.ReadFile(m)
	if err != nil {
		return "This destination host is not trusted", 500
	}
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            rootCAs,
	}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{
		Transport: tr,
		Timeout:   500 * time.Millisecond,
	}
	// end of trusting ca
	payload := Payload{"shit"}
	encodepayload, _ := json.Marshal(payload)
	ebody := bytes.NewReader(encodepayload)
	url := "https://test2.klin-pro.com:2018/"
	req, err := http.NewRequest("PATCH", url, ebody)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-API-Key", "wtf")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err.Error(), 500
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body), 200
}

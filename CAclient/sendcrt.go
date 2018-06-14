package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, csr []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(csr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, err
}

func sendcsr(m, host, k string, csr []byte) ([]byte, error) {
	var toReturn []byte
	request, err := newfileUploadRequest(host, csr)
	if err != nil {
		return toReturn, err
	}

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	certs, err := ioutil.ReadFile(m)
	if err != nil {
		return toReturn, err
	}
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		return toReturn, errors.New("No certs appended, using system certs only")
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            rootCAs,
	}
	config.BuildNameToCertificate()
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{
		Transport: tr,
		Timeout:   1500 * time.Millisecond,
	}
	//** end of clean up
	resp, err := client.Do(request)
	if err != nil {
		return toReturn, err
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		toReturn = body.Bytes()
	}
	return toReturn, nil
}

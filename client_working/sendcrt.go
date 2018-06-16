package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}

func sendcsrv2(mastercert, host string) {
	extraParams := map[string]string{
		"filename": "testfilename",
	}
	request, err := newfileUploadRequest(host, extraParams, "file", "testfile")
	if err != nil {
		panic(err)
	}

	cert, err := tls.LoadX509KeyPair("test2.klin-pro.com.crt", "test2.klin-pro.com.key")
	if err != nil {
		log.Fatalln("Unable to load cert", err)
	}

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	certs, err := ioutil.ReadFile(mastercert)
	if err != nil {
		log.Fatal(err)
	}
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            rootCAs,
	}
	config.BuildNameToCertificate()
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{
		Transport: tr,
		Timeout:   500 * time.Millisecond,
	}
	//** end of clean up
	fmt.Println("wtf")
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode)
}

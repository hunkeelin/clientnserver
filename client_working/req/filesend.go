package klinreq

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func customRequest(uri, paramName, path, method string, params map[string]string) (*http.Request, error) {
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

	req, err := http.NewRequest(method, uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err

}

// Creates a new file upload http request with optional extra params
func Sendfile(i *ReqInfo) {
	extraParams := map[string]string{
		"filename": "klinFile",
	}
	req, err := customRequest("https://"+i.Dest+":"+i.Dport, "file", i.File, i.Method, extraParams)
	if err != nil {
		panic(err)
	}
	cert, err := tls.LoadX509KeyPair(i.Cert, i.Key)
	if err != nil {
		log.Fatalln("Unable to load cert", err)
	}

	// Load our CA certificate
	clientCACert, err := ioutil.ReadFile(i.Trust)
	if err != nil {
		log.Fatal("Unable to open cert", err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		//	InsecureSkipVerify: false,
		Certificates: []tls.Certificate{cert},
		RootCAs:      clientCertPool,
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{
		Timeout:   500 * time.Millisecond,
		Transport: tr,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body), string(resp.Status))
}

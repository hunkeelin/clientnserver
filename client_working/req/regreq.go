package klinreq

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
func SendPayload(i *ReqInfo, payload interface{}) {
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
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            clientCertPool,
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{
		Timeout:   500 * time.Millisecond,
		Transport: tr,
	}
	//	payload := Payload{C: "asdf", D: true}

	encodepayload, _ := json.Marshal(payload)
	ebody := bytes.NewReader(encodepayload)
	req, err := http.NewRequest(i.Method, "https://"+i.Dest+":"+i.Dport, ebody)
	if err != nil {
		log.Println("Unable to speak to our server", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body), string(resp.Status))
}

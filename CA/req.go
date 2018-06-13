package main

import (
	"crypto/x509"
	"fmt"
	"github.com/hunkeelin/pki"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.Header.Get("content-type"), "application/x-www-form-urlencoded") {
		d, err := ioutil.ReadAll(r.Body)
		rawcert, err := klinpki.SignCSRv2("program/CA/ca.crt", "program/CA/ca.key", d, 300)
		if err != nil {
			panic(err)
		}
		clientCSR, _ := x509.ParseCertificateRequest(d)
		hostname, _ := net.LookupAddr(strings.Split(r.RemoteAddr, ":")[0])
		if clientCSR.DNSNames[0]+"." == hostname[0] {
			w.WriteHeader(202)
			w.Write(rawcert)
		} else {
			fmt.Println("fuck off")
		}
	}
	return
}

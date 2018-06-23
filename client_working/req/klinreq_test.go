package klinreq

import (
	"fmt"
	"testing"
)

type testPayload struct {
	C string `json:"content"`
	D bool   `json:"disabled"`
}

func TestReq(t *testing.T) {
	fmt.Println("testing req")
	i := &ReqInfo{
		Cert:   "../test2.klin-pro.com.crt",
		Key:    "../test2.klin-pro.com.key",
		Dest:   "test3.klin-pro.com",
		Dport:  "2018",
		Trust:  "../master.pem",
		Method: "POST",
	}
	payload := &testPayload{
		C: "wtf",
		D: true,
	}
	SendPayload(i, payload)
}

func TestSendfile(t *testing.T) {
	fmt.Println("testing filesend")
	i := &ReqInfo{
		Cert:   "../test2.klin-pro.com.crt",
		Key:    "../test2.klin-pro.com.key",
		Dest:   "test3.klin-pro.com",
		Dport:  "2018",
		Trust:  "../master.pem",
		Method: "POST",
		File:   "../testfile",
	}
	Sendfile(i)
}

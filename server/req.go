package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %v! \n", r.TLS.PeerCertificates[0].EmailAddresses[0])
	if strings.HasPrefix(r.Header.Get("content-type"), "multipart/form-data") {
		f.monorun <- struct{}{} //makine sure there's no two request writing to the same file
		msg, status := mgmtcert(r, f.pkidir)
		<-f.monorun
		w.WriteHeader(status)
		w.Write([]byte(msg))
		fmt.Println(msg, status)
		return
	} else {
		return
	}
}

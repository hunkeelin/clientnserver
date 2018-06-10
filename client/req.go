package main

import (
	"fmt"
	"net/http"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header.Get("content-type"))
	status := 200
	msg := "it worked"
	w.WriteHeader(status)
	w.Write([]byte(msg))
	fmt.Println(msg, status)
	return
}

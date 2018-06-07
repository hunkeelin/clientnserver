package main

import (
	"fmt"
	"net/http"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
	status := 200
	msg := "it worked"
	w.WriteHeader(status)
	w.Write([]byte(msg))
	fmt.Println(msg, status)
	return
}

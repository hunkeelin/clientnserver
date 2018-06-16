package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.Header.Get("content-type"), "multipart/form-data") {
		t, _, _ := r.FormFile("file")
		to, _ := os.Create("shit")
		io.Copy(to, t)
		to.Close()
		msg := "fag"
		status := 100
		fmt.Println(msg, status)
		w.WriteHeader(status)
		w.Write([]byte(msg))
		return
	} else {
		msg := "wut"
		status := 200
		fmt.Println(msg, status)
		w.WriteHeader(status)
		w.Write([]byte(msg))
		return
	}
}

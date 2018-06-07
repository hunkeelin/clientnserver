package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func mgmtcert(r *http.Request, pkidir string) (string, int) {
	t, _, err := r.FormFile("file")
	if err != nil {
		return "unable to form file", 500
	}
	crtname := r.FormValue("filename")
	crtdir := pkidir + "pending/" + crtname
	if Exist(crtdir) {
		return "crt already exist", 500
	}
	to, err := os.Create(crtdir)
	if err != nil {
		return "unable to create file", 500
	}
	io.Copy(to, t)
	err = to.Close()
	if err != nil {
		return "unable to close file " + crtdir, 500
	}
	g := r.FormValue("filename")
	fmt.Println("wrote ", g, " to pending")
	return "it worked", 200
}

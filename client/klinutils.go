package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
)

func dowork() {
	fmt.Println("start")
	time.Sleep(9 * time.Second)
	fmt.Println("end")
}
func joblist(path string) map[string]string {
	m := make(map[string]string)
	err := os.Chdir(path)
	if err != nil {
		log.Fatal("no such file or directory ", path)
	}
	ls, _ := filepath.Glob("*")
	for _, provider := range ls { //github.com and other provider
		os.Chdir(path + provider)
		orgs, _ := filepath.Glob("*")
		for _, org := range orgs { //github.com/orgs
			os.Chdir(path + provider + "/" + org)
			jobs, _ := filepath.Glob("*")
			for _, job := range jobs { // going through each jobs
				confdir, _ := os.Getwd()
				confdir = confdir + "/" + job + "/" + "config.xml"
				url := "https://" + provider + "/" + org + "/" + job
				m[strings.ToLower(url)] = confdir

			}
		}
	}
	return m
}
func getHostnameFromCert(path string) string {
	e, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("can't read file")
	}
	block, _ := pem.Decode(e)
	if block == nil {
		log.Fatal("fail to parse certifiate pem")
	}
	leaf, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("can't parse file", path)
		log.Fatal(err)
	}
	fmt.Println(len(leaf.DNSNames))
	return leaf.DNSNames[0]
}

func isvalidmethod(r *http.Request) bool {
	methodlist := []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTOINS", "TRACE"}
	return stringInSlice(r.Method, methodlist)
}

func matchstring(s, regex string) bool {
	match, err := regexp.MatchString(regex, s)
	fmt.Println("")
	if err != nil {
		log.Fatal("regex matching problem")
	}
	return match
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func exist(p string) bool {
	if _, err := os.Stat(p); err == nil {
		return true
	}
	return false
}
func runshell(cmd string, args []string, uid, gid uint32) error {
	//err := exec.Command(cmd, args...).Run()
	acmd := exec.Command(cmd, args...)
	acmd.SysProcAttr = &syscall.SysProcAttr{}
	acmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}
	out, err := acmd.CombinedOutput()
	fmt.Printf("%s\n", out)
	return err
}

func removestring(s []string, pattern string) []string {
	var toreturn []string
	for _, raw_element := range s {
		element := strings.Replace(raw_element, " ", "", -1)
		if strings.HasPrefix(element, pattern) {
			ele := strings.Replace(element, pattern, "", -1)
			if strings.HasPrefix(ele, "HEAD") || strings.HasPrefix(ele, "master") {
				continue
			}
			toreturn = append(toreturn, ele)
		}
	}
	return toreturn
}

func outshell(cmd string, args []string) (string, error) {
	output, err := exec.Command(cmd, args...).Output()
	return string(output), err
}
func cleandir(s, env []string, workers int) {
	sema := make(chan struct{}, workers)
	wg := sync.WaitGroup{}
	for _, element := range s {
		if !stringInSlice(element, env) {
			sema <- struct{}{}
			wg.Add(1)
			go func(element string) {
				os.RemoveAll(element)
				<-sema
				wg.Done()
			}(element)
		}
	}
	wg.Wait()
}

func createdir(s, env []string, workers int) {
	sema := make(chan struct{}, workers)
	wg := sync.WaitGroup{}
	for _, element := range env {
		if stringInSlice(element, s) == false {
			sema <- struct{}{}
			wg.Add(1)
			go func(element string) {
				os.MkdirAll(element+"/"+"modules", 0755)
				<-sema
				wg.Done()
			}(element)
		}
	}
	wg.Wait()
}
func checkstring(s, pattern string) {
	if strings.HasPrefix(s, "mod") {
		if string(s[len(s)-1]) != "," {
			log.Fatal("error missing comma on line: ", s)
		}
	}

	checknum := 0

	for _, r := range s {
		c := string(r)
		if c == "'" {
			checknum = checknum + 1
		}
	}

	if checknum != 2 {
		log.Fatal("error missing single quotes on line: ", s)
	}
}
func checkerr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func trim(x string) string {
	pattern := "'"
	checkstring(x, pattern)
	bra := strings.Index(x, pattern)
	if bra < 0 {
		return ""
	}
	rx := x[bra+1:]
	ket := strings.Index(rx, pattern)
	return rx[:ket]
}

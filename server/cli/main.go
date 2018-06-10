package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

var (
	hosts = flag.String("nodes", "*", "regex of the hosts avalible")
)

func main() {
	flag.Parse()
	pkidir := "/home/bgops/files/golesson/clientnserver/server/program/"
	files, _ := filepath.Glob(pkidir + "trust/" + *hosts)
	for _, f := range files {
		fmt.Println(f)
	}
	msg, status := dowork(pkidir + "trust/" + "test2.klin-pro.com.crt")
	fmt.Println(msg, status)
}

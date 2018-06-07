package main

import (
	"fmt"
)

func main() {
	pkidir := "/home/bgops/files/golesson/clientnserver/server/program/"
	msg, status := dowork(pkidir + "trust/" + "test2.klin-pro.com.crt")
	fmt.Println(msg, status)
}

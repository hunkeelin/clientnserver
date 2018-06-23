package main

import ()

func main() {
	testreq("master.pem", "http://test3.klin-pro.com:2018")
	sendfile("master.pem", "https://test3.klin-pro.com:2018")
}

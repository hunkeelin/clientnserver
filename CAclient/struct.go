package main

import (
	"sync"
)

type Config struct {
	apikey     string
	bindaddr   string
	pkidir     string
	port       string
	masterport string
	keycertdir string
	masteraddr string
	concur     int
	jobdir     string
	mastercrt  string
}
type Conn struct {
	regex   string
	apikey  string
	pkidir  string
	concur  int
	jobdir  string
	mu      sync.Mutex
	monorun chan struct{}
}

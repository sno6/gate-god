package main

import (
	"github.com/sno6/gate-god/server/ftp"
)

func main() {
	server := ftp.New(nil)
	server.Serve()
}

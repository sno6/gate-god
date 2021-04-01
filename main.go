package main

import (
	"log"

	"github.com/sno6/gate-god/camera"
	"github.com/sno6/gate-god/server/ftp"
)

func main() {
	batcher := camera.NewFrameBatcher()
	server := ftp.New(&ftp.Config{
		User:     "admin",
		Password: "password",
	}, batcher)

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"log"

	cmd "github.com/sno6/gate-god/cmd/gate-god"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

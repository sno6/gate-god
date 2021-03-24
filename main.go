package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sno6/gate-god/recognition/platerecognizer"
)

func main() {
	// server := ftp.New(nil)
	// server.Serve()

	f, err := os.Open("./image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := platerecognizer.New()
	result, err := r.Recognize(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

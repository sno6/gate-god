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

	token := os.Getenv("PLATE_RECOGNIZER_API_TOKEN")
	if token == "" {
		panic("empty token")
	}

	r := platerecognizer.New(token)
	result, err := r.RecognizePlate(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

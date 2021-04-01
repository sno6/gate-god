package camera

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

const tickerTimeout = time.Second * 3

type Frame struct {
	R    io.Reader
	Name string
}

// FrameBatcher is a service that runs concurrently
// and stitches together frames into batches as they come in
// from the server.
//
// The purpose of the FrameBatcher is to be able to select the best
// frame out of a batch, as well as for folder level storage.
type FrameBatcher struct {
	frameChan chan *Frame
	logger    *log.Logger

	bufferT *time.Ticker
	buffer  []*Frame
}

func NewFrameBatcher() *FrameBatcher {
	fb := &FrameBatcher{
		frameChan: make(chan *Frame),
		logger:    log.New(os.Stdout, "[Frame]: ", log.LstdFlags),
		bufferT:   time.NewTicker(tickerTimeout),
	}
	go fb.process()
	return fb
}

func (fb *FrameBatcher) SendFrame(f *Frame) {
	go func() {
		fb.frameChan <- f
	}()
}

func (fb *FrameBatcher) process() {
	for {
		select {
		case <-fb.bufferT.C:
			if len(fb.buffer) == 0 {
				fmt.Println("No motion")
				break
			}

			// Save copy of buffer to in separate go routine to free processing.
			fmt.Printf("Saving %d frames in batch\n", len(fb.buffer))
			go saveBuffer(fb.buffer[:])

			// Reset the buffer for the next batch.
			fb.buffer = nil
		case f := <-fb.frameChan:
			fb.logger.Printf("Received frame: %s\n", f.Name)

			// Unfortunately, because of some crazy issue within the FTP server library
			// we are using we need to drain the reader before the handler can accept
			// a new PUT file request. More info: https://gitea.com/goftp/server/issues/148
			r, err := readerToReader(f.R)
			if err != nil {
				log.Printf("readerToReader: %v\n", err)
			}

			fb.buffer = append(fb.buffer, &Frame{R: r, Name: f.Name})
			fb.bufferT.Reset(tickerTimeout)
		}
	}
}

func readerToReader(r io.Reader) (io.Reader, error) {
	b := &bytes.Buffer{}
	_, err := io.Copy(b, r)
	return b, err
}

// Testing purposes.. probably wanna chuck this in S3 later..
func saveBuffer(frames []*Frame) {
	dirName := path.Join("./motion", strconv.Itoa(int(time.Now().UnixNano())))
	if err := os.Mkdir(dirName, 0777); err != nil {
		panic(err)
	}

	for _, frame := range frames {
		f, err := os.Create(path.Join(dirName, frame.Name))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = io.Copy(f, frame.R)
		if err != nil {
			panic(err)
		}
	}
}

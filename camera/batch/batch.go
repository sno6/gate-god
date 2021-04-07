package batch

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/sno6/gate-god/camera"
	"github.com/sno6/gate-god/engine"
)

const tickerTimeout = time.Second * 3

// FrameBatcher is a service that runs concurrently
// and stitches together frames into batches as they come in
// from the server (on motion detection).
//
// The purpose of the FrameBatcher is to be able to select the best
// frame(s) out of a batch, as well as for storage.
type FrameBatcher struct {
	engine *engine.Engine

	frameChan chan *camera.Frame
	logger    *log.Logger

	bufferT *time.Ticker
	buffer  []*camera.Frame
}

func New(engine *engine.Engine) *FrameBatcher {
	fb := &FrameBatcher{
		engine:    engine,
		frameChan: make(chan *camera.Frame),
		logger:    log.New(os.Stdout, "[Frame]: ", log.LstdFlags),
		bufferT:   time.NewTicker(tickerTimeout),
	}
	go fb.process()
	return fb
}

func (fb *FrameBatcher) SendFrame(f *camera.Frame) {
	go func() {
		fb.frameChan <- f
	}()
}

func (fb *FrameBatcher) process() {
	for {
		select {
		case <-fb.bufferT.C:
			if len(fb.buffer) == 0 {
				break
			}

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

			frame := &camera.Frame{R: r, Name: f.Name}

			// If we are in "key frame" territory, meaning we are around 3-4 frames
			// in, start sending them to the engine to process license plates and
			// core logic. This is done to skip "car approaching" frames.
			if fb.isKeyFrame(len(fb.buffer)) {
				fb.engine.SendFrame(frame)
			}

			fb.buffer = append(fb.buffer, frame)
			fb.bufferT.Reset(tickerTimeout)
		}
	}
}

// We determine a key frame as being the 3rd -> 6th frame in the shot.
func (fb *FrameBatcher) isKeyFrame(bufLen int) bool {
	return bufLen >= 3 && bufLen <= 6
}

func readerToReader(r io.Reader) (io.Reader, error) {
	b := &bytes.Buffer{}
	_, err := io.Copy(b, r)
	return b, err
}

// Testing purposes.. probably wanna chuck this in S3 later..
func saveBuffer(frames []*camera.Frame) {
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
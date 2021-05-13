package batch

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"

	"github.com/sno6/gate-god/camera"
	"github.com/sno6/gate-god/engine"
)

const frameSaveDir = "./motion"
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

	currentBatchID int
	bufferT        *time.Ticker
	buffer         []*camera.Frame
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

			// Reset the buffer for the next batch.
			fb.currentBatchID++
			fb.buffer = nil
		case f := <-fb.frameChan:
			fb.logger.Printf("Received frame: %s\n", f.Name)

			// Unfortunately, because of some crazy issue within the FTP server library
			// we are using we need to drain the reader before the handler can accept
			// a new PUT file request. More info: https://gitea.com/goftp/server/issues/148
			r, err := readerToReader(f.Reader)
			if err != nil {
				log.Printf("error copying reader: %v\n", err)
				break
			}

			// If we are in "key frame" territory, meaning we are around 1-4 frames
			// in, start sending them to the engine to process license plates and
			// core logic. This is done to skip "car approaching" frames.
			if fb.isKeyFrame(len(fb.buffer)) {
				fb.engine.SendFrame(&camera.Frame{
					Name:    f.Name,
					BatchID: fb.currentBatchID,
					Reader:  r,
				})
			}

			fb.buffer = append(fb.buffer, &camera.Frame{
				Name:    f.Name,
				BatchID: fb.currentBatchID,
			})

			fb.bufferT.Reset(tickerTimeout)
		}
	}
}

// We determine a key frame as being the 1st -> 4th frame in the shot.
func (fb *FrameBatcher) isKeyFrame(bufLen int) bool {
	return bufLen >= 1 && bufLen <= 4
}

func readerToReader(r io.Reader) (io.Reader, error) {
	b := &bytes.Buffer{}
	_, err := io.Copy(b, r)
	return b, err
}

package engine

import (
	"log"
	"os"
	"strings"

	"github.com/sno6/gate-god/camera"
	"github.com/sno6/gate-god/recognition"
	"github.com/sno6/gate-god/relay"
)

type Engine struct {
	recognizer    recognition.PlaterRecognizer
	relay         *relay.Relay
	allowedPlates []string

	frameChan          chan *camera.Frame
	currentBatchID     int
	ignoreCurrentBatch bool

	logger *log.Logger
}

func New(
	recognizer recognition.PlaterRecognizer,
	relay *relay.Relay,
	allowedPlated []string,
) *Engine {
	e := &Engine{
		recognizer: recognizer,
		relay:      relay,
		frameChan:  make(chan *camera.Frame),
		logger:     log.New(os.Stdout, "[Engine]: ", log.LstdFlags),
	}
	go e.process()
	return e
}

func (e *Engine) SendFrame(f *camera.Frame) {
	go func() {
		e.frameChan <- f
	}()
}

func (e *Engine) process() {
	for {
		select {
		case f := <-e.frameChan:
			// If we see a new batch, reset state.
			if f.BatchID != e.currentBatchID {
				e.currentBatchID = f.BatchID
				e.ignoreCurrentBatch = false
			}

			// Already found what we wanted, ignore new frames.
			if e.ignoreCurrentBatch {
				break
			}

			e.logger.Printf("Engine received prime frame: %s\n", f.Name)

			// Run the plate through recognition.
			plate, err := e.recognizer.RecognizePlate(f.Reader)
			if err != nil {
				e.logger.Println(err)
				break
			}

			e.logger.Printf("Received plate: %v with score: %v\n", plate.Plate, plate.Score)

			if !e.isPlateAllowed(plate.Plate) {
				e.logger.Printf("Access attempt rejected.")
				break
			} else {
				// We don't need to process anymore images in this batch.
				e.ignoreCurrentBatch = true
			}

			// The gate is now open, enjoy your stay.
			e.relay.Toggle()
		}
	}
}

func (e *Engine) isPlateAllowed(plate string) bool {
	for _, p := range e.allowedPlates {
		if strings.ToLower(p) == strings.ToLower(plate) {
			return true
		}
	}

	return false
}

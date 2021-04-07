package engine

import (
	"fmt"

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

			fmt.Printf("Engine received prime frame: %s\n", f.Name)

			// TODO:
			// 1. Run through recognition.
			// 2. Check allowed list for plate.
			// 3. Open the gate.
		}
	}
}

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

	frameChan chan *camera.Frame
}

func New(
	recognizer recognition.PlaterRecognizer,
	relay *relay.Relay,
	allowedPlated []string,
) *Engine {
	e := &Engine{
		recognizer: recognizer,
		relay:      relay,
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
			fmt.Printf("Engine received prime frame: %s\n", f.Name)

			// 1. Run through recognition.
			// 2. Check allowed list for plate.
			// 3. Open the gate.
		}
	}
}

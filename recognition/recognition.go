package recognition

import "io"

type Result struct {
}

type Recognizer interface {
	RecognizePlate(r io.Reader) (*Result, error)
}

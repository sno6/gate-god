package recognition

import "io"

type Result struct {
	Score float64 `json:"score"`
	Plate string  `json:"plate"`
}

type Recognizer interface {
	RecognizePlate(r io.Reader) (*Result, error)
}

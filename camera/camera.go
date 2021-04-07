package camera

import "io"

type Frame struct {
	Reader  io.Reader
	Name    string
	BatchID int
}

package camera

import "io"

type Frame struct {
	R    io.Reader
	Name string
}

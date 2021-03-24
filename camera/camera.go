package camera

import "io"

type MotionDetector interface {
	OnMotionDetection(r io.Reader) error
}

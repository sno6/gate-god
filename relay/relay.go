package relay

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

type Relay struct {
	pin rpio.Pin
}

func New(mcuPin int) (*Relay, error) {
	r := &Relay{
		pin: rpio.Pin(mcuPin),
	}
	err := r.setup()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func NewDummy(mcuPin int) (*Relay, error) {
	return &Relay{pin: rpio.Pin(mcuPin)}, nil
}

func (r *Relay) setup() error {
	defer r.pin.Output()
	return rpio.Open()
}

func (r *Relay) Toggle() {
	r.pin.High()
	time.Sleep(time.Second)
	r.pin.Low()
}

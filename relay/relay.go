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

func (r *Relay) setup() error {
	return rpio.Open()
}

func (r *Relay) Toggle() {
	r.pin.Toggle()
}

func (r *Relay) TestRelay() {
	for i := 0; i < 5; i++ {
		r.Toggle()

		time.Sleep(time.Second * 10)
	}
}

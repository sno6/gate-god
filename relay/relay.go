package relay

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

type Relayer interface {
	Toggle()
}

type Relay struct {
	pin rpio.Pin
}

type DummyRelay struct{}

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

func NewDummy(_ int) (*DummyRelay, error) {
	return &DummyRelay{}, nil
}

func (d *DummyRelay) Toggle() {}

func (r *Relay) setup() error {
	defer r.pin.Output()
	return rpio.Open()
}

func (r *Relay) Toggle() {
	r.pin.High()
	time.Sleep(time.Second * 3)
	r.pin.Low()
}

func (r *Relay) TestRelay() {
	r.Toggle()
	time.Sleep(time.Second * 10)
	r.Toggle()
}

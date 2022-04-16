package gpioconfig

import (
	"github.com/stianeikeland/go-rpio"
)

type pin struct {
	allowed bool
	rpioPin rpio.Pin
}

type GPIOPortRPI4 struct {
	pins map[uint8]pin
	open bool
}

func (g GPIOPortRPI4) IsConfigured(gpioPin uint8) bool {
	_, err := g.pins[gpioPin]
	return err
}

func (g GPIOPortRPI4) IsAllowed(gpioPin uint8) bool {
	return g.pins[gpioPin].allowed
}

// Requires g.IsOpen()
func (g *GPIOPortRPI4) SetPinAsCoil(gpioPin uint8) error {
	p := rpio.Pin(gpioPin)
	p.Output()
	g.pins[gpioPin] = pin{allowed: false, rpioPin: p}
	return nil
}

func (g *GPIOPortRPI4) Allow(gpioPin uint8) error {
	pin := g.pins[gpioPin]
	pin.allowed = true
	g.pins[gpioPin] = pin
	return nil
}

func (g *GPIOPortRPI4) Deny(gpioPin uint8) error {
	pin := g.pins[gpioPin]
	pin.allowed = false
	g.pins[gpioPin] = pin
	return nil
}

func (g *GPIOPortRPI4) Open() error {
	if err := rpio.Open(); err != nil {
		return err
	}
	g.pins = make(map[uint8]pin)
	g.open = true
	return nil
}

func (g *GPIOPortRPI4) Close() error {
	if err := rpio.Close(); err != nil {
		return err
	}
	g.open = false
	return nil
}

func (g GPIOPortRPI4) IsOpen() bool {
	return g.open
}

// Requires g.IsOpen()
func (g GPIOPortRPI4) SetCoil(gpioPin uint8, val bool) {
	if val {
		g.pins[gpioPin].rpioPin.High()
	} else {
		g.pins[gpioPin].rpioPin.Low()
	}
}

func (g GPIOPortRPI4) GetCoil(gpioPin uint8) (res bool) {
	if g.pins[gpioPin].rpioPin.Read() == rpio.Low {
		return false
	} else {
		return true
	}
}

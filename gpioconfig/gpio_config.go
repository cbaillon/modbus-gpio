package gpioconfig

import (
	"github.com/stianeikeland/go-rpio"
)

type pin struct {
	allowed bool
	rpioPin rpio.Pin
}

type GPIOPort struct {
	pins map[uint8]pin
	open bool
}

func (g GPIOPort) IsConfigured(gpioPin uint8) bool {
	_, err := g.pins[gpioPin]
	return err
}

func (g GPIOPort) IsAllowed(gpioPin uint8) bool {
	return g.pins[gpioPin].allowed
}

// Requires g.IsOpen()
func (g *GPIOPort) SetPinAsCoil(gpioPin uint8) error {
	p := rpio.Pin(gpioPin)
	p.Output()
	g.pins[gpioPin] = pin{allowed: false, rpioPin: p}
	return nil
}

func (g *GPIOPort) Allow(gpioPin uint8) error {
	pin := g.pins[gpioPin]
	pin.allowed = true
	g.pins[gpioPin] = pin
	return nil
}

func (g *GPIOPort) Deny(gpioPin uint8) error {
	pin := g.pins[gpioPin]
	pin.allowed = false
	g.pins[gpioPin] = pin
	return nil
}

func (g *GPIOPort) Open() error {
	if err := rpio.Open(); err != nil {
		return err
	}
	g.pins = make(map[uint8]pin)
	g.open = true
	return nil
}

func (g *GPIOPort) Close() error {
	if err := rpio.Close(); err != nil {
		return err
	}
	g.open = false
	return nil
}

func (g GPIOPort) IsOpen() bool {
	return g.open
}

// Requires g.IsOpen()
func (g GPIOPort) SetCoil(gpioPin uint8, val bool) {
	if val {
		g.pins[gpioPin].rpioPin.High()
	} else {
		g.pins[gpioPin].rpioPin.Low()
	}
}

func (g GPIOPort) GetCoil(GPIOPort uint8) (res bool) {
	if g.pins[GPIOPort].rpioPin.Read() == rpio.Low {
		return false
	} else {
		return true
	}
}

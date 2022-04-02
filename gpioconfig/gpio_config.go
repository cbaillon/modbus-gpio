package gpioconfig

import (
	"fmt"

	"github.com/stianeikeland/go-rpio"
)

type pin struct {
	allowed bool
	rpioPin rpio.Pin
}

type gpioPort struct {
	pins map[uint8]pin
	open bool
}

func (g gpioPort) IsConfigured(gpioPin uint8) bool {
	_, err := g.pins[gpioPin]
	return err
}

func (g gpioPort) IsAllowed(gpioPin uint8) bool {
	return g.pins[gpioPin].allowed
}

// Requires g.IsOpen()
func (g *gpioPort) SetPinAsCoil(gpioPin uint8) error {
	p := rpio.Pin(gpioPin)
	p.Output()
	fmt.Println("g: ", g, "g.pins:", g.pins)
	//g.pins[gpioPin] = pin{allowed: true, rpioPin: p}
	g.pins[gpioPin] = pin{allowed: true, rpioPin: p}
	return nil
}

func (g *gpioPort) Open() error {
	if err := rpio.Open(); err != nil {
		return err
	}
	g.pins = make(map[uint8]pin)
	g.open = true
	return nil
}

func (g *gpioPort) Close() error {
	if err := rpio.Close(); err != nil {
		return err
	}
	g.open = false
	return nil
}

func (g gpioPort) IsOpen() bool {
	return g.open
}

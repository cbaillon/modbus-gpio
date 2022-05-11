package gpioconfig

import (
	"errors"
	"fmt"

	"github.com/stianeikeland/go-rpio"
)

type pin struct {
	allowed  bool
	rpioPin  rpio.Pin
	pullMode PullMode
}

type GPIOPort struct {
	pins map[uint8]pin
	open bool
}

type PullMode uint8

// Pull Up / Down
const (
	PullModeOff PullMode = iota
	PullModeDown
	PullModeUp
)

func (g GPIOPort) IsConfigured(gpioPin uint8) bool {
	_, err := g.pins[gpioPin]
	return err
}

func (g GPIOPort) IsAllowed(gpioPin uint8) bool {
	fmt.Println("Call of IsAllowed")
	g.PrintPortConfiguration()
	return g.pins[gpioPin].allowed
}

func (g GPIOPort) PrintPortConfiguration() {
	fmt.Println("***** Port Configuration")
	fmt.Println("Port open: ", g.open)
	for k, v := range g.pins {
        fmt.Printf("%s -> %s\n", k, v)
    }
	fmt.Println("************************")
	
}

// Requires g.IsOpen()
func (g *GPIOPort) SetPinAsCoil(gpioPin uint8) error {
	p := rpio.Pin(gpioPin)
	p.Output()
	g.pins[gpioPin] = pin{allowed: false, rpioPin: p}
	return nil
}

// Requires g.IsOpen()
// To set a pin as a Discrete Input (in the sense of Modbus terminology), we have
// to specify the PullMode of the pin. The mode arg must be PullModeDown or PullModeUp.
// If you don't know how GPIO inputs behavior with pull-up or pull-down modes work,
// you can find out more with the following article:
// https://kalitut.com/raspberrypi-gpio-pull-up-pull-down-resistor/
func (g *GPIOPort) SetPinAsDiscreteInput(gpioPin uint8, mode PullMode) error {
	if mode != PullModeDown && mode != PullModeUp {
		return errors.New("invalid mode")
	}

	p := rpio.Pin(gpioPin)
	p.Input()
	if mode == PullModeDown {
		p.PullDown()
	} else if mode == PullModeUp {
		p.PullUp()
	}
	g.pins[gpioPin] = pin{allowed: false, rpioPin: p, pullMode: mode}
	return nil
}

// Requires g.IsOpen() && gpioPin must be configured as a DiscreteInput
func (g GPIOPort) GetPinPullMode(gpioPin uint8) PullMode {
	return g.pins[gpioPin].pullMode
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

func (g GPIOPort) GetCoil(gpioPort uint8) (res bool) {
	if g.pins[gpioPort].rpioPin.Read() == rpio.Low {
		return false
	} else {
		return true
	}
}

func (g GPIOPort) GetDiscreteInput(gpioPort uint8) (bool, error) {
	val := g.pins[gpioPort].rpioPin.Read()
	pullmode := g.GetPinPullMode(gpioPort)

	switch pullmode {
	case PullModeDown:
		switch val {
		case rpio.Low:
			return false, nil
		case rpio.High:
			return true, nil
		default:
			return false, errors.New("bad result of rpio.Read(): " + string(val))
		}
	case PullModeUp:
		switch val {
		case rpio.Low:
			return true, nil
		case rpio.High:
			return false, nil
		default:
			return false, errors.New("bad result of rpio.Read(): " + string(val))
		}
	default:
		return false, errors.New("wrong PullMode:" + string(pullmode))
	}
}


func MapDiscreteInputToInputRegister(gpioPort uint8, inputRegisterAddress uint16, bitoffset, uint8){

}

func MapCoilToHoldingRegister(gpioPort uint8, HoldingRegisterAddress uint16, bitoffset, uint8){

}
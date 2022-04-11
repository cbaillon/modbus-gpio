package gpioconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsConfiguredAndIsAllowed(t *testing.T) {
	var gp GPIOPort

	// Make sure alls pins are not configured when
	// not initialised.
	want := false
	for a := uint8(0); a <= 40; a++ {
		if ic := gp.IsConfigured(a); ic != want {
			t.Errorf("IsConfigured: With a=%d, Got %t, Want %t", a, ic, want)
		}
		if ia := gp.IsAllowed(a); ia != want {
			t.Errorf("IsAllowed: With a=%d, Got %t, Want %t", a, ia, want)

		}
	}
}

func TestIsOpen(t *testing.T) {
	var gp GPIOPort
	if gp.IsOpen() {
		t.Error("gpioPort should not be open by default")
	}

	err := gp.Open()
	if err != nil {
		t.Errorf("Error while opening: %s", err)
		return
	}
	if !gp.IsOpen() {
		t.Error("gpioPort should be open after Open() call")
	}

	gp.Close()
	if gp.IsOpen() {
		t.Error("gpioPort should not be open after calling Close()")
	}

}

func TestSetPinAsCoil(t *testing.T) {
	var gp GPIOPort
	err := gp.Open()
	if err != nil {
		t.Errorf("Error while opening: %s", err)
		return
	}
	assert.Equal(t, len(gp.pins), 0, "length of pins map should be 0 just after opening.")
	gp.SetPinAsCoil(17)
	assert.Equal(t, len(gp.pins), 1, "length of pins map should 1 after calling SetPinAsCoil once.")
	assert.Equal(t, int(gp.pins[17].rpioPin), 17, "pin at index x should be x (here: 17)")
}

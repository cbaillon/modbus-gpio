package gpioconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultIsConfiguredAndIsAllowed(t *testing.T) {
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

func TestIsAllowed(t *testing.T) {
	var gp GPIOPort
	err := gp.Open()
	if err != nil {
		t.Errorf("Error while opening: %s", err)
		return
	}

	assert.Equal(t, false, gp.IsAllowed(22), "pin 22 should not be allowed before configuring")
	gp.Allow(22)
	assert.Equal(t, true, gp.IsAllowed(22), "pin 22 should be allowed after configuring")

	assert.Equal(t, false, gp.IsAllowed(27), "pin 27 should not be allowed before configuring")
	gp.Allow(27)
	assert.Equal(t, true, gp.IsAllowed(27), "pin 27 should be allowed after configuring")
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

	assert.Equal(t, 0, len(gp.pins), "length of pins map should be 0 just after opening.")
	assert.Equal(t, false, gp.IsConfigured(17), "pin should not be configured by default")
	assert.Equal(t, false, gp.IsAllowed(17), "pin should not be allowed by default")

	gp.SetPinAsCoil(17)
	assert.Equal(t, 1, len(gp.pins), "length of pins map should be 1 after calling SetPinAsCoil once.")
	assert.Equal(t, 17, int(gp.pins[17].rpioPin), "pin at index x should be x (here: 17)")
	assert.Equal(t, true, gp.IsConfigured(17), "pin should be configured after call of SetPinAsCoil")
	assert.Equal(t, false, gp.IsAllowed(17), "pin should not be allowed by default if Allow() not explicitly called")

	gp.Allow(17)
	assert.Equal(t, true, gp.IsConfigured(17), "pin should still be configured after call of Allow")
	assert.Equal(t, true, gp.IsAllowed(17), "pin should be allowed after call of Allow()")

	gp.Deny(17)
	assert.Equal(t, true, gp.IsConfigured(17), "pin should still be configured after call of Deny")
	assert.Equal(t, false, gp.IsAllowed(17), "pin should be denied after call of Deny()")

	gp.Close()
}

func TestCannotAllowUnconfiguredPin(t *testing.T) {
	var gp GPIOPort
	err := gp.Open()
	if err != nil {
		t.Errorf("Error while opening: %s", err)
		return
	}
	// to be done
}

package gpioconfig

type GPIOPort interface {

	// Requires: IsOpen() == true
	// Is the pin given in argument is configured on the hardware
	// Note that a configured pin means it can be read or write locally, but doesn't imply
	// that the pin is allowed to be accessed from a Modbus request.
	IsConfigured(gpioPin uint8) bool

	// Requires: IsOpen() == true
	// Is the pin is allowed to be accessed from a Modbus request.
	IsAllowed(gpioPin uint8) bool

	// Requires: IsOpen() == true
	// Configure a pin on the GPIO port as an output.
	// It will be mapped as a Coil in the Modbus terminology.
	// Coils in Modbus are read-write single bits.
	SetPinAsCoil(gpioPin uint8) error

	// Requires: IsOpen() == true
	// Requires: IsConfigured(gpioPin) == true
	// Allow a previously configured pin of the GPIO port
	// to be accessed via a Modbus request.
	*Allow(gpioPin uint8) error

	// Requires: IsOpen() == true
	// Requires: IsConfigured(gpioPin) == true
	// Deny a previously configured pin of the GPIO port
	// to be accessed via a Modbus request.
	Deny(gpioPin uint8) error

	// Open the GPIO port.
	Open() error

	// Close the GPIO port.
	Close() error

	// Is the GPIO port open?
	IsOpen() bool

	// Requires: IsOpen() == true
	// Set the value of a pin that is configured as Coil.
	SetCoil(gpioPin uint8, val bool)

	// Requires: IsOpen() == true
	// Return the value of a pin that is configured as Coil.
	GetCoil(GPIOPort uint8) (res bool)
}

func MakeGPIOPortRPI4() GPIOPort {
	return GPIOPortRPI4{}
}


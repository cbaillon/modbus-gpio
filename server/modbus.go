package server

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cbaillon/modbus-gpio/gpioconfig"
	"github.com/simonvetter/modbus"
)

const (
	MINUS_ONE int16 = -1
)

/*
* Start_Server simply start listening for Modbus requests.
* This function never ends
 */

func Start_server(gpioPort *gpioconfig.GPIOPort, ip string, port string) {
	var server *modbus.ModbusServer
	var err error
	var eh *modbusGPIOHandler = &modbusGPIOHandler{}

	eh.port = gpioPort

	// create the server object
	server, err = modbus.NewServer(&modbus.ServerConfiguration{
		URL: "tcp://" + ip + ":" + port,
		// close idle connections after 30s of inactivity
		Timeout: 30 * time.Second,
		// accept 15 concurrent connections max.
		MaxClients: 15,
	}, eh)
	if err != nil {
		fmt.Printf("failed to create server: %v\n", err)
		os.Exit(1)
	}

	// start accepting client connections
	// note that Start() returns as soon as the server is started
	err = server.Start()
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(1)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}

// Modbus handler object
type modbusGPIOHandler struct {
	// this lock is used to avoid concurrency issues between goroutines, as
	// handler methods are called from different goroutines
	// (1 goroutine per client)
	lock sync.RWMutex

	port *gpioconfig.GPIOPort
}

// Requires bitofsset >= 0 < 16
func setBitHighAt(n *uint8, bitoffset uint8) {
	*n = *n | (uint8(2) ^ bitoffset)
}

// Requires bitofsset >= 0 < 16
func setBitLowAt(n *uint8, bitoffset uint8) {
	*n = *n&uint8(2) ^ bitoffset
}

// IntPow calculates n to the mth power. Since the result is an int, it is assumed that m is a positive power
// Go doesn't natively support integer exponents!!
func intPow(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}

// Coil handler method.
// This method gets called whenever a valid modbus request asking for a coil operation is
// received by the server.
func (mh *modbusGPIOHandler) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	var addr uint8 = uint8(req.Addr)
	fmt.Println("Starting HandleCoils. Addr = ", addr)
	fmt.Println("UnitId = ", req.UnitId)
	if req.UnitId != 255 {
		// only accept unit ID #255
		err = modbus.ErrBadUnitId
		return nil, errors.New("HandleCoils: error, UnitId must be 255, was " + string(req.UnitId))
	}

	if req.Quantity != 1 {
		err = modbus.ErrIllegalDataValue
		return nil, errors.New("HandleCoils: error, only requests with quantity of 1 allowed" + fmt.Sprint(req.Quantity))
	}

	if !mh.port.IsAllowed(addr) {
		fmt.Println("addr", addr, " is not allowed. Returning an error")
		err = modbus.ErrIllegalDataAddress
		return nil, errors.New(string(modbus.ErrIllegalDataAddress) + " - HandleCoils: error, coils at address" + string(addr) + " is not allowed")
	} else {
		fmt.Println("Coil at addr ", addr, "is allowed")
	}

	if req.IsWrite {
		// since we're manipulating variables shared between multiple goroutines,
		// acquire a lock to avoid concurrency issues.
		mh.lock.Lock()
		// release the lock upon return
		defer mh.lock.Unlock()
		fmt.Println("req.Args: ", req.Args)
		mh.port.SetCoil(addr, req.Args[0])
	} else {
		res = append(res, mh.port.GetCoil(addr))
	}
	return
}

// Discrete input handler method.
// Note that we're returning ErrIllegalFunction unconditionally.
// This will cause the client to receive "illegal function", which is the modbus way of
// reporting that this server does not support/implement the discrete input type.
func (mh *modbusGPIOHandler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	var addr uint8 = uint8(req.Addr)
	if req.UnitId != 255 {
		// only accept unit ID #255
		// note: we're merely filtering here, but we could as well use the unit
		// ID field to support multiple register maps in a single server.
		err = modbus.ErrBadUnitId
		return nil, errors.New("HandleCoilds: error, UnitId must be 255, was " + string(req.UnitId))
	}

	if req.Quantity != 1 {
		err = modbus.ErrIllegalDataValue
		return nil, errors.New("HandleCoilds: error, only requests with quantity of 1 allowed" + fmt.Sprint(req.Quantity))
	}

	if !mh.port.IsAllowed(addr) {
		err = modbus.ErrIllegalDataAddress
		return nil, errors.New(string(modbus.ErrIllegalDataAddress) + " - HandleCoils: error, coils at address" + string(addr) + " is not allowed")
	}

	val, err := mh.port.GetDiscreteInput(addr)

	if err != nil {
		return nil, err
	}
	res = append(res, val)
	return
}

// Holding register handler method.
// Not supported. Return ErrIllegalFunction error
func (eh *modbusGPIOHandler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	var regAddr uint16

	if req.UnitId != 1 {
		// only accept unit ID #1
		err = modbus.ErrIllegalFunction
		return
	}

	// since we're manipulating variables shared between multiple goroutines,
	// acquire a lock to avoid concurrency issues.
	eh.lock.Lock()
	// release the lock upon return
	defer eh.lock.Unlock()

	// loop through `quantity` registers
	for i := 0; i < int(req.Quantity); i++ {
		// compute the target register address
		regAddr = req.Addr + uint16(i)

		switch regAddr {
		// expose the static, read-only value of 0xff00 in register 100
		case 100:
			res = append(res, 0xff00)

		// expose holdingReg1 in register 101 (RW)
		case 101:
			if req.IsWrite {
				eh.holdingReg1 = req.Args[i]
			}
			res = append(res, eh.holdingReg1)

		// expose holdingReg2 in register 102 (RW)
		case 102:
			if req.IsWrite {
				// only accept values 2 and 4
				switch req.Args[i] {
				case 2, 4:
					eh.holdingReg2 = req.Args[i]

					// make note of the change (e.g. for auditing purposes)
					fmt.Printf("%s set reg#102 to %v\n", req.ClientAddr, eh.holdingReg2)
				default:
					// if the written value is neither 2 nor 4,
					// return a modbus "illegal data value" to
					// let the client know that the value is
					// not acceptable.
					err = modbus.ErrIllegalDataValue
					return
				}
			}
			res = append(res, eh.holdingReg2)

		// expose eh.holdingReg3 in register 103 (RW)
		// note: eh.holdingReg3 is a signed 16-bit integer
		case 103:
			if req.IsWrite {
				// cast the 16-bit unsigned integer passed by the server
				// to a 16-bit signed integer when writing
				eh.holdingReg3 = int16(req.Args[i])
			}
			// cast the 16-bit signed integer from the handler to a 16-bit unsigned
			// integer so that we can append it to `res`.
			res = append(res, uint16(eh.holdingReg3))

		// expose the 16 most-significant bits of eh.holdingReg4 in register 200
		case 200:
			if req.IsWrite {
				eh.holdingReg4 =
					((uint32(req.Args[i])<<16)&0xffff0000 |
						(eh.holdingReg4 & 0x0000ffff))
			}
			res = append(res, uint16((eh.holdingReg4>>16)&0x0000ffff))

		// expose the 16 least-significant bits of eh.holdingReg4 in register 201
		case 201:
			if req.IsWrite {
				eh.holdingReg4 =
					(uint32(req.Args[i])&0x0000ffff |
						(eh.holdingReg4 & 0xffff0000))
			}
			res = append(res, uint16(eh.holdingReg4&0x0000ffff))

		// any other address is unknown
		default:
			err = modbus.ErrIllegalDataAddress
			return
		}
	}

	return
}

// Input register handler method.
// Not supported. Return ErrIllegalFunction error
func (eh *modbusGPIOHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	err = modbus.ErrIllegalFunction
	return
}

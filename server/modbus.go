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
		return nil, errors.New("HandleCoilds: error, UnitId must be 255, was " + string(req.UnitId))
	}

	if req.Quantity != 1 {
		err = modbus.ErrIllegalDataValue
		return nil, errors.New("HandleCoilds: error, only requests with quantity of 1 allowed" + fmt.Sprint(req.Quantity))
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
	err = modbus.ErrIllegalFunction
	return
}

// Input register handler method.
// Not supported. Return ErrIllegalFunction error
func (eh *modbusGPIOHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	err = modbus.ErrIllegalFunction
	return
}

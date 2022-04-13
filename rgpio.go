package main

import (
	"fmt"
	"os"

	"github.com/cbaillon/modbus-gpio/gpioconfig"
	"github.com/cbaillon/modbus-gpio/server"
)

func main() {
	fmt.Println("rgpio starting")

	port := gpioconfig.GPIOPort{}

	if err := port.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer port.Close()

	port.SetPinAsCoil(17)
	port.Allow(17)

	server.Start_server(&port)
}

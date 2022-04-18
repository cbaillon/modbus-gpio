package main

import (
	"fmt"
	"os"
	"time"

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

	port.SetPinAsDiscreteInput(6, gpioconfig.PullModeUp)
	port.Allow(6)

	go func() {
		time.Sleep(time.Second * 2)
		for {
			pin6, err := port.GetDiscreteInput(6)
			if err != nil {
				fmt.Println("Error while reading DiscreteInput")
				os.Exit(1)
			}
			if pin6 {
				fmt.Println("Button is pressed")
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	server.Start_server(&port)
}

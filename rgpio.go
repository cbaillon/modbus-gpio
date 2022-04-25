package main

import (
	"fmt"
	"os"
	//"time"
	

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

	fmt.Println("Port configuration before setting:")
	port.PrintPortConfiguration()

	port.SetPinAsDiscreteInput(17, gpioconfig.PullModeDown)
	port.Allow(17)

	port.SetPinAsCoil(27)
	port.Allow(27)

	port.SetPinAsCoil(22)
	port.Allow(22)

	fmt.Println("Port configuration after setting:")
	port.PrintPortConfiguration()
	
	/* 	go func() {
		time.Sleep(time.Second * 2)
		for {
			pin17, err := port.GetDiscreteInput(17)
			if err != nil {
				fmt.Println("Error while reading DiscreteInput")
				os.Exit(1)
			}
			if pin17 {
				//fmt.Println("Contact is on")
			} else {
				//fmt.Println("Contact is off")
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()


	b := true
	for {
		port.SetCoil(22, b)
		fmt.Println("Pin 22 set to ", b)
		
		time.Sleep(5*time.Second)

		port.SetCoil(27, b)
		fmt.Println("Pin 27 set to ", b)

		b = ! b
		time.Sleep(5*time.Second)
	} */


	server.Start_server(&port, "192.168.0.81", "5502")
}

IP=192.168.0.84
PORT=5502


build:
	go build rgpio.go

run:
	./rgpio $(IP) $(PORT)

create_gpio_permissions: # Create permissions for the user launching this Makefile
	sudo python3 create_gpio_user_permissions.py $(USER)
	sudo udevadm trigger

lock_door:
	modbus-cli --target tcp://$(IP):$(PORT) -unit-id 255 wc:27:true

unlock_door:
	modbus-cli --target tcp://$(IP):$(PORT) -unit-id 255 wc:27:false

open_strike:
	modbus-cli --target tcp://$(IP):$(PORT) -unit-id 255 wc:27:true
	sleep 2
	modbus-cli --target tcp://$(IP):$(PORT) -unit-id 255 wc:27:false

get_door_state:
	modbus-cli --target tcp://$(IP):$(PORT) -unit-id 255 rdi:17


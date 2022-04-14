


build:
	go build rgpio.go

create_gpio_permissions: # Create permissions for the user launching this Makefile
	sudo python3 create_gpio_user_permissions.py $(USER)
	sudo udevadm trigger


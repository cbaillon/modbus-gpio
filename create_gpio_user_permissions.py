# This script comes from: https://github.com/wuestkamp/raspberry-gpio-python/blob/master/create_gpio_user_permissions.py

import grp
import sys
import subprocess


def ensure_gpiogroup():
    try:
        grp.getgrnam('gpio')
        print('GPIO group alread exists')
    except KeyError:
        print('GPIO group does not exist - creating...')
        subprocess.call(['groupadd', '-f', '-r', 'gpio'])
       
    
def add_user_to_gpio_group(user):
    subprocess.call(['adduser', user, 'gpio'])
    # in future, also for groups:
    #   spi
    #   i2c

def add_udev_rules():
    with open('/etc/udev/rules.d/99-gpio.rules','w') as f:
        f.write("""SUBSYSTEM=="bcm2835-gpiomem", KERNEL=="gpiomem", GROUP="gpio", MODE="0660"
SUBSYSTEM=="gpio", KERNEL=="gpiochip*", ACTION=="add", PROGRAM="/bin/sh -c 'chown root:gpio /sys/class/gpio/export /sys/class/gpio/unexport ; chmod 220 /sys/class/gpio/export /sys/class/gpio/unexport'"
SUBSYSTEM=="gpio", KERNEL=="gpio*", ACTION=="add", PROGRAM="/bin/sh -c 'chown root:gpio /sys%p/active_low /sys%p/direction /sys%p/edge /sys%p/value ; chmod 660 /sys%p/active_low /sys%p/direction /sys%p/edge /sys%p/value'"
""")

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print('Syntax: create_gpio_user_permissions.py USER_NAME')
        sys.exit(1)

    username = sys.argv[1]
    print("Creating GPIO permissions for user ", username, "...")
    ensure_gpiogroup()
    add_user_to_gpio_group(username)
    add_udev_rules()


I2C-bus interaction of peripheral sensors with Raspberry PI embedded linux or respective clones
==============================================================================================

[![Build Status](https://travis-ci.org/d2r2/go-i2c.svg?branch=master)](https://travis-ci.org/d2r2/go-i2c)
[![Go Report Card](https://goreportcard.com/badge/github.com/d2r2/go-i2c)](https://goreportcard.com/report/github.com/d2r2/go-i2c)
[![GoDoc](https://godoc.org/github.com/d2r2/go-i2c?status.svg)](https://godoc.org/github.com/d2r2/go-i2c)
[![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

This library written in [Go programming language](https://golang.org/) intended to activate and interact with the I2C bus by reading and writing data.

Compatibility
-------------

Tested on Raspberry Pi 1 (model B), Raspberry Pi 3 B+, Banana Pi (model M1), Orange Pi Zero, Orange Pi One.

Golang usage
------------

```go
func main() {
  // Create new connection to I2C bus on 2 line with address 0x27
  i2c, err := i2c.NewI2C(0x27, 2)
  if err != nil { log.Fatal(err) }
  // Free I2C connection on exit
  defer i2c.Close()
  ....
  // Here goes code specific for sending and reading data
  // to and from device connected via I2C bus, like:
  _, err := i2c.Write([]byte{0x1, 0xF3})
  if err != nil { log.Fatal(err) }
  ....
}
```

Tutorial
--------

My [repositories](https://github.com/d2r2?tab=repositories) contain quite a lot projects, which use i2c library as a starting point to interact with various peripheral devices and sensors for use on embedded Linux devices. All these libraries start with a standard call to open I2C-connection to specific bus line and address, than pass i2c instance to device.

In its turn, go-i2c use [go-logger](https://github.com/d2r2/go-logger) library to output debug and other notification's lines which produce all necessary levels of logging. You can manage what level of verbosity you would like to see, by adding call:
```go
// Uncomment/comment next line to suppress/increase verbosity of output
logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
```
Once you put this call, it will decrease verbosity from default "Debug" up to next "Info" level, reducing the number of low-level console outputs that occur during interaction with the I2C bus. Please, find examples in corresponding I2C-driven sensors among my projects.

You will find here the list of all devices and sensors supported by me, that reference this library:

- [Liquid-crystal display driven by Hitachi HD44780 IC](https://github.com/d2r2/go-hd44780).
- [BMP180/BMP280/BME280 temperature and pressure sensors](https://github.com/d2r2/go-bsbmp).
- [DHT12/AM2320 humidity and temperature sensors](https://github.com/d2r2/go-aosong).
- [Si7021 relative humidity and temperature sensor](https://github.com/d2r2/go-si7021).
- [SHT3x humidity and temperature sensor](https://github.com/d2r2/go-sht3x).
- [VL53L0X time-of-flight ranging sensor](https://github.com/d2r2/go-vl53l0x).
- [BH1750 ambient light sensor](https://github.com/d2r2/go-bh1750).
- [MPL3115A2 pressure and temperature sensor](https://github.com/d2r2/go-mpl3115a2).


Getting help
------------

GoDoc [documentation](http://godoc.org/github.com/d2r2/go-i2c)

Troubleshooting
--------------

- *How to obtain fresh Golang installation to RPi device (either any RPi clone):*
If your RaspberryPI golang installation taken by default from repository is outdated, you may consider
to install actual golang manually from official Golang [site](https://golang.org/dl/). Download
tar.gz file containing armv6l in the name. Follow installation instructions.

- *How to enable I2C bus on RPi device:*
If you employ RaspberryPI, use raspi-config utility to activate i2c-bus on the OS level.
Go to "Interfacing Options" menu, to active I2C bus.
Probably you will need to reboot to load i2c kernel module.
Finally you should have device like /dev/i2c-1 present in the system.

- *How to find I2C bus allocation and device address:*
Use i2cdetect utility in format "i2cdetect -y X", where X may vary from 0 to 5 or more,
to discover address occupied by peripheral device. To install utility you should run
`apt install i2c-tools` on debian-kind system. `i2cdetect -y 1` sample output:
	```
	     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
	00:          -- -- -- -- -- -- -- -- -- -- -- -- --
	10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	70: -- -- -- -- -- -- 76 --    
	```

License
-------

Go-i2c is licensed under MIT License.

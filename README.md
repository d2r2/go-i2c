## I2C bus setting up and usage for linux on RPi device and respective clones

This library written in [Go programming language](https://golang.org/) intended to activate and interact with the I2C bus by reading and writing data.

## Compatibility

Tested on Raspberry PI 1 (model B) and Banana PI (model M1).

## Golang usage

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

## Getting help

GoDoc [documentation](http://godoc.org/github.com/d2r2/go-i2c)

## Troubleshoting

- How to obtain fresh Golang installation to RPi device (either any RPi clone):
  
  Download fresh stable ARM tar.gz release file (containing armv6l in file name): https://golang.org/dl/.
  Read instruction how to unpack content to /usr/local/ folder and update/set up such variables from user environment as PATH, GOPATH and so on.

- How to enable I2C bus on RPi device:
  
  Your /dev/ folder should contains files like /dev/i2c-1 to have i2c support activated in the kernel. Otherwise you should find proper module to active it via `modprobe` utility, either config it permanently via /etc/modules config file.

- How to find display I2C bus and address:

  Use i2cdetect utility in format "i2cdetect -y X", where X vary from 0 to 5 or more, to discover address occupied by device. To install utility you should run `apt-get install i2c-tools` on debian-kind system.

## License

Go-i2c is licensed inder MIT License.

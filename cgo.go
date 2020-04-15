// +build linux,cgo

package i2c

// #include <linux/i2c-dev.h>
import "C"

// Get I2C_SLAVE constant value from
// Linux OS I2C declaration file.
var (
	I2C_SLAVE = GetValI2c()
)

func GetValI2c() uintptr {
	if C.I2C_SLAVE != 0 {
		return C.I2C_SLAVE
	}
	return 0x0703
}

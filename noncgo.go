// +build !cgo

package i2c

// Use hard-coded value for system I2C_SLAVE
// constant, if OS not Linux or CGO disabled.
// This is not a good approach, but
// can be used as a last resort.
const (
	I2C_SLAVE       = 0x0703
	I2C_SLAVE_FORCE = 0x0706
)

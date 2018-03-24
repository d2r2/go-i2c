// Package i2c provides low level control over the linux i2c bus.
//
// Before usage you should load the i2c-dev kernel module
//
//      sudo modprobe i2c-dev
//
// Each i2c bus can address 127 independent i2c devices, and most
// linux systems contain several buses.
package i2c

import (
	"fmt"
	"log"
	"os"
	"syscall"

	logger "github.com/d2r2/go-logger"
)

var lg = logger.NewPackageLogger("i2c",
	// logger.DebugLevel,
	logger.InfoLevel,
)

const (
	I2C_SLAVE = 0x0703
)

// I2C represents a connection to an i2c device.
type I2C struct {
	rc *os.File
	// Logger
	log *log.Logger
	// Enable verbose output
	Debug bool
}

// New opens a connection to an i2c device.
func NewI2C(addr uint8, bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := ioctl(f.Fd(), I2C_SLAVE, uintptr(addr)); err != nil {
		return nil, err
	}
	this := &I2C{rc: f}
	return this, nil
}

// Write sends buf to the remote i2c device. The interpretation of
// the message is implementation dependant.
func (this *I2C) write(buf []byte) (int, error) {
	return this.rc.Write(buf)
}

func (this *I2C) WriteByte(b byte) (int, error) {
	var buf [1]byte
	buf[0] = b
	return this.rc.Write(buf[:])
}

func (this *I2C) read(p []byte) (int, error) {
	return this.rc.Read(p)
}

func (this *I2C) Close() error {
	return this.rc.Close()
}

// SMBus (System Management Bus) protocol over I2C.
// Read count of n byte's sequence from i2c device
// starting from reg address.
func (this *I2C) ReadRegBytes(reg byte, n int) ([]byte, int, error) {
	_, err := this.write([]byte{reg})
	if err != nil {
		return nil, 0, err
	}
	buf := make([]byte, n)
	c, err := this.read(buf)
	if err != nil {
		return nil, 0, err
	}
	lg.Debugf("Read %d bytes starting from reg 0x%0X", c, reg)
	return buf, c, nil

}

// SMBus (System Management Bus) protocol over I2C.
// Read byte from i2c device register specified in reg.
func (this *I2C) ReadRegU8(reg byte) (byte, error) {
	_, err := this.write([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 1)
	_, err = this.read(buf)
	if err != nil {
		return 0, err
	}
	lg.Debugf("Read U8 %d from reg 0x%0X", buf[0], reg)
	return buf[0], nil
}

// SMBus (System Management Bus) protocol over I2C.
// Write byte to i2c device register specified in reg.
func (this *I2C) WriteRegU8(reg byte, value byte) error {
	buf := []byte{reg, value}
	_, err := this.write(buf)
	if err != nil {
		return err
	}
	lg.Debugf("Write U8 %d to reg 0x%0X", value, reg)
	return nil
}

// SMBus (System Management Bus) protocol over I2C.
// Read unsigned big endian word (16 bits) from i2c device
// starting from address specified in reg.
func (this *I2C) ReadRegU16BE(reg byte) (uint16, error) {
	_, err := this.write([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = this.read(buf)
	if err != nil {
		return 0, err
	}
	w := uint16(buf[0])<<8 + uint16(buf[1])
	lg.Debugf("Read U16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// SMBus (System Management Bus) protocol over I2C.
// Read unsigned little endian word (16 bits) from i2c device
// starting from address specified in reg.
func (this *I2C) ReadRegU16LE(reg byte) (uint16, error) {
	w, err := this.ReadRegU16BE(reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil
}

// SMBus (System Management Bus) protocol over I2C.
// Read signed big endian word (16 bits) from i2c device
// starting from address specified in reg.
func (this *I2C) ReadRegS16BE(reg byte) (int16, error) {
	_, err := this.write([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = this.read(buf)
	if err != nil {
		return 0, err
	}
	w := int16(buf[0])<<8 + int16(buf[1])
	lg.Debugf("Read S16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// SMBus (System Management Bus) protocol over I2C.
// Read unsigned little endian word (16 bits) from i2c device
// starting from address specified in reg.
func (this *I2C) ReadRegS16LE(reg byte) (int16, error) {
	w, err := this.ReadRegS16BE(reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil

}

// SMBus (System Management Bus) protocol over I2C.
// Write unsigned big endian word (16 bits) value to i2c device
// starting from address specified in reg.
func (this *I2C) WriteRegU16BE(reg byte, value uint16) error {
	buf := []byte{reg, byte((value & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := this.write(buf)
	if err != nil {
		return err
	}
	lg.Debugf("Write U16 %d to reg 0x%0X", value, reg)
	return nil
}

// SMBus (System Management Bus) protocol over I2C.
// Write unsigned big endian word (16 bits) value to i2c device
// starting from address specified in reg.
func (this *I2C) WriteRegU16LE(reg byte, value uint16) error {
	w := (value*0xFF00)>>8 + value<<8
	return this.WriteRegU16BE(reg, w)
}

// SMBus (System Management Bus) protocol over I2C.
// Write signed big endian word (16 bits) value to i2c device
// starting from address specified in reg.
func (this *I2C) WriteRegS16BE(reg byte, value int16) error {
	buf := []byte{reg, byte((uint16(value) & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := this.write(buf)
	if err != nil {
		return err
	}
	lg.Debugf("Write S16 %d to reg 0x%0X", value, reg)
	return nil
}

// SMBus (System Management Bus) protocol over I2C.
// Write signed big endian word (16 bits) value to i2c device
// starting from address specified in reg.
func (this *I2C) WriteRegS16LE(reg byte, value int16) error {
	w := int16((uint16(value)*0xFF00)>>8) + value<<8
	return this.WriteRegS16BE(reg, w)
}

func ioctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}

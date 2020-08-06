// Package i2c provides low level control over the Linux i2c bus.
//
// Before usage you should load the i2c-dev kernel module
//
//      sudo modprobe i2c-dev
//
// Each i2c bus can address 127 independent i2c devices, and most
// Linux systems contain several buses.
package i2c

import (
	"encoding/hex"
	"fmt"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Options represents a connection to I2C-device.
type Options struct {
	addr uint8
	bus  int
	rc   *os.File
	Log  *logrus.Logger
}

// New opens a connection for I2C-device.
// SMBus (System Management Bus) protocol over I2C
// supported as well: you should preliminary specify
// register address to read from, either write register
// together with the data in case of write operations.
func New(addr uint8, bus int) (*Options, error) {
	v := &Options{
		addr: addr,
		bus:  bus,
		Log: &logrus.Logger{
			Out:       os.Stderr,
			Formatter: new(logrus.TextFormatter),
			//Hooks:     make(logrus.LevelHooks),
			Level: logrus.InfoLevel,
		},
	}

	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return v, err
	}
	if err := ioctl(f.Fd(), I2C_SLAVE, uintptr(addr)); err != nil {
		return v, err
	}

	v.rc = f
	return v, nil
}

// GetBus return bus line, where I2C-device is allocated.
func (o *Options) GetBus() int {
	return o.bus
}

// GetAddr return device occupied address in the bus.
func (o *Options) GetAddr() uint8 {
	return o.addr
}

func (o *Options) write(buf []byte) (int, error) {
	return o.rc.Write(buf)
}

// WriteBytes send bytes to the remote I2C-device. The interpretation of
// the message is implementation-dependent.
func (o *Options) WriteBytes(buf []byte) (int, error) {
	o.Log.Debugf("Write %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return o.write(buf)
}

func (o *Options) read(buf []byte) (int, error) {
	return o.rc.Read(buf)
}

// ReadBytes read bytes from I2C-device.
// Number of bytes read correspond to buf parameter length.
func (o *Options) ReadBytes(buf []byte) (int, error) {
	n, err := o.read(buf)
	if err != nil {
		return n, err
	}
	o.Log.Debugf("Read %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return n, nil
}

// Close I2C-connection.
func (o *Options) Close() error {
	return o.rc.Close()
}

// ReadRegBytes read count of n byte's sequence from I2C-device
// starting from reg address.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) ReadRegBytes(reg byte, n int) ([]byte, int, error) {
	o.Log.Debugf("Read %d bytes starting from reg 0x%0X...", n, reg)
	_, err := o.WriteBytes([]byte{reg})
	if err != nil {

		return nil, 0, err
	}
	buf := make([]byte, n)
	c, err := o.ReadBytes(buf)
	if err != nil {
		return nil, 0, err
	}
	return buf, c, nil

}

// ReadRegU8 reads byte from I2C-device register specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) ReadRegU8(reg byte) (byte, error) {
	_, err := o.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 1)
	_, err = o.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	o.Log.Debugf("Read U8 %d from reg 0x%0X", buf[0], reg)
	return buf[0], nil
}

// WriteRegU8 writes byte to I2C-device register specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) WriteRegU8(reg byte, value byte) error {
	buf := []byte{reg, value}
	_, err := o.WriteBytes(buf)
	if err != nil {
		return err
	}
	o.Log.Debugf("Write U8 %d to reg 0x%0X", value, reg)
	return nil
}

// ReadRegU16BE reads unsigned big endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) ReadRegU16BE(reg byte) (uint16, error) {
	_, err := o.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = o.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	w := uint16(buf[0])<<8 + uint16(buf[1])
	o.Log.Debugf("Read U16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// ReadRegU16LE reads unsigned little endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) ReadRegU16LE(reg byte) (uint16, error) {
	w, err := o.ReadRegU16BE(reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil
}

// ReadRegS16BE reads signed big endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) ReadRegS16BE(reg byte) (int16, error) {
	_, err := o.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = o.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	w := int16(buf[0])<<8 + int16(buf[1])
	o.Log.Debugf("Read S16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// ReadRegS16LE reads signed little endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) ReadRegS16LE(reg byte) (int16, error) {
	w, err := o.ReadRegS16BE(reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil

}

// WriteRegU16BE writes unsigned big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) WriteRegU16BE(reg byte, value uint16) error {
	buf := []byte{reg, byte((value & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := o.WriteBytes(buf)
	if err != nil {
		return err
	}
	o.Log.Debugf("Write U16 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegU16LE writes unsigned little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) WriteRegU16LE(reg byte, value uint16) error {
	w := (value*0xFF00)>>8 + value<<8
	return o.WriteRegU16BE(reg, w)
}

// WriteRegS16BE writes signed big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) WriteRegS16BE(reg byte, value int16) error {
	buf := []byte{reg, byte((uint16(value) & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := o.WriteBytes(buf)
	if err != nil {
		return err
	}
	o.Log.Debugf("Write S16 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegS16LE writes signed little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (o *Options) WriteRegS16LE(reg byte, value int16) error {
	w := int16((uint16(value)*0xFF00)>>8) + value<<8
	return o.WriteRegS16BE(reg, w)
}

func ioctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}

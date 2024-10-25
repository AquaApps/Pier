package tunnel

import (
	"context"
	"io"

	//"log"
	"net"
	"os"
)

const MTU = 1500

type Device struct {
	Name   string
	CIDRv4 *net.IPNet

	_f          *os.File
	_life       context.Context
	_cancelFunc context.CancelFunc

	_interceptFunc func([]byte) []byte
}

func (device *Device) Init(mName, mCIDRv4 string, ctx context.Context) error {
	device._life, device._cancelFunc = context.WithCancel(ctx)
	device.Name = mName
	if IPv4, CIDRv4, err := net.ParseCIDR(mCIDRv4); err != nil {
		return err
	} else {
		device.CIDRv4 = CIDRv4
		device.CIDRv4.IP = IPv4
	}

	if f, err := openTunDeviceWithIP(device.Name, device.CIDRv4, MTU); err != nil {
		return err
	} else {
		device._f = f
	}
	device._interceptFunc = func(bytes []byte) []byte {
		//log.Println(bytes)
		return bytes
	}
	return nil
}

func (device *Device) OpenChannel() (io.Reader, io.Writer) {
	return device._f, device._f
}

func (device *Device) Destroy() {
	device._cancelFunc()
	close(device._inputStream)
	close(device._outputStream)
	closeTunDevice(device._f)
}

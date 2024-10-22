package tunnel

import (
	"context"
	//"log"
	"net"
	"os"
	"sync/atomic"
)

const MTU = 1500

type Device struct {
	Name   string
	CIDRv4 *net.IPNet

	_inputStream  chan []byte
	_outputStream chan []byte

	_f                 *os.File
	_life              context.Context
	_cancelFunc        context.CancelFunc
	_totalReadBytes    *uint64
	_totalWrittenBytes *uint64
	_interceptFunc     func([]byte) []byte
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

	device._totalReadBytes = new(uint64)
	device._totalWrittenBytes = new(uint64)
	device._inputStream = make(chan []byte)
	device._outputStream = make(chan []byte)
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

func (device *Device) OpenChannel() (out <-chan []byte, in chan<- []byte) {
	go device.readFromTunnel()
	go device.writeToTunnel()
	return device._outputStream, device._inputStream
}

func (device *Device) GetReadWriteBytes() (uint64, uint64) {
	return atomic.LoadUint64(device._totalReadBytes), atomic.LoadUint64(device._totalWrittenBytes)
}

func (device *Device) Destroy() {
	device._cancelFunc()
	closeTunDevice(device._f)
}

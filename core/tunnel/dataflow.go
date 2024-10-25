package tunnel

import (
	"Pier/common"
	"sync/atomic"
)

func (device *Device) readFromTunnel() {
	fun := device._interceptFunc
	packet := make([]byte, 64*1024)
	for common.Opened(device._life) {
		num, _ := device._f.Read(packet)
		device.incrWrittenBytes(num)
		device._outputStream <- fun(packet[:num])
	}
}

func (device *Device) writeToTunnel() {
	fun := device._interceptFunc
	for common.Opened(device._life) {
		n, _ := device._f.Write(fun(<-device._inputStream))
		device.incrReadBytes(n)
	}
}

func (device *Device) incrReadBytes(n int) {
	atomic.AddUint64(device._totalReadBytes, uint64(n))
}

func (device *Device) incrWrittenBytes(n int) {
	atomic.AddUint64(device._totalWrittenBytes, uint64(n))
}

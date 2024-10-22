package tunnel

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
	"os"
	"syscall"
	"unsafe"
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func openTunDeviceWithIP(name string, v4 *net.IPNet, mtu int) (*os.File, error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	var req ifReq
	copy(req.Name[:], name)
	req.Flags = syscall.IFF_TUN | syscall.IFF_NO_PI
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		_ = file.Close()
		err = fmt.Errorf("ioctl %v", errno)
		return nil, err
	}

	netInterface, err := netlink.LinkByName(name)

	if err != nil {
		return nil, err
	}

	addrV4 := &netlink.Addr{IPNet: v4, Label: ""}

	if err = netlink.LinkSetMTU(netInterface, mtu); err != nil {
		return nil, err
	}

	if err = netlink.AddrAdd(netInterface, addrV4); err != nil {
		return nil, err
	}

	if err = netlink.LinkSetUp(netInterface); err != nil {
		return nil, err
	}
	return file, nil
}

func closeTunDevice(f *os.File) {
	_ = f.Close()
}

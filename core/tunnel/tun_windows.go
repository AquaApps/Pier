package tunnel

import (
	"github.com/net-byte/water"
	"net"
	"os"
)

func openTunDeviceWithIP(name string, v4 *net.IPNet, mtu int) (*os.File, error) {
	// todo:impl https://github.com/getlantern/gotun/blob/master/gotun_windows.go
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name:    name,
			Network: []string{v4.String()},
		},
	})
	if err != nil {
		return nil, err
	}
	return ifce, nil
}

func closeTunDevice(f *os.File) {
	_ = f.Close()
}

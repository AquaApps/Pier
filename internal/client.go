package internal

import (
	"Pier/core/transport"
	"Pier/core/tunnel"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
)

func PierClient(Addr string, CIDRv4 string, dev *tunnel.Device, appContext context.Context, transportImpl transport.Transport) error {
	if err := injectRouteForClient(CIDRv4); err != nil {
		log.Println("[Client]Warning:", err)
	}
	defer unInjectRouteForClient(CIDRv4)
	out, in := dev.OpenChannel()
	return transportImpl.Dail(appContext, Addr, func(stream io.ReadWriter) {
		log.Println("Connected to service.")
		connectCtx, cancelFunc := context.WithCancel(appContext)
		defer cancelFunc()
		go writer(stream, out, connectCtx)
		reader(stream, in, connectCtx)
	})
}

func injectRouteForClient(CIDRv4 string) error {
	IPv4, NETv4, err := net.ParseCIDR(CIDRv4)
	if err != nil {
		return err
	}
	commands := [][]string{
		{"ip", "route", "add", NETv4.String(), "via", IPv4.String()},
	}

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			return fmt.Errorf("failed to execute command '%x %s': %w", cmd[0], cmd[1:], err)
		}
	}
	return nil
}

func unInjectRouteForClient(CIDRv4 string) {
	_, NETv4, _ := net.ParseCIDR(CIDRv4)
	commands := [][]string{
		{"ip", "route", "delete", NETv4.String()},
	}
	for _, cmd := range commands {
		_ = exec.Command(cmd[0], cmd[1:]...).Run()
	}
}

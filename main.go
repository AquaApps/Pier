package main

import (
	"Pier/common"
	"Pier/core/service"
	"Pier/core/transport"
	"Pier/internal"
	"context"
	pun "github.com/AquaApps/Pun"
	"log"
	"net"
)

var appContext = context.Background()
var transportImpl = new(transport.TcpTransport)
var punDev *pun.Device

func main() {
	if err := CreatePun(); err != nil {
		log.Fatal("[Initial]Fail to init tunDevice:", err)
		return
	}

	log.Println("[Initial]Init tunDevice success!")
	defer punDev.Close()

	if config.HttpService.Enable {
		go func() {
			log.Printf("[HttpService]Start listen %d\n", config.HttpService.Port)
			if err := service.StartHttpServer(config.HttpService.Port); err != nil {
				log.Printf("[HttpService]Listen %d fail: %v\n", config.HttpService.Port, err)
			}
		}()
	}
	if config.ServerMode {
		if err := internal.PierServer(config.ServiceAddr, punDev, appContext, transportImpl); err != nil {
			log.Fatal("[Server]Fail to start Service.", err)
			return
		}
	} else {
		if err := internal.PierClient(config.ServiceAddr, config.CIDRv4, punDev, appContext, transportImpl); err != nil {
			log.Fatal("[Client]Fail to start Client.", err)
			return
		}
	}
	log.Println("Unknown State. Waiting signal...")
	common.WaitingSignal()
}

func CreatePun() error {
	var devName string
	if config.Extra.ObfName {
		devName = common.ObfuscateText(config.TunName)
	} else {
		devName = config.TunName
	}
	var mCIDRv4 net.IPNet
	if IPv4, CIDRv4, err := net.ParseCIDR(config.CIDRv4); err != nil {
		return err
	} else {
		mCIDRv4 = *CIDRv4
		mCIDRv4.IP = IPv4
	}
	if dev, err := pun.New(&pun.Config{
		Name:   devName,
		CIDRv4: mCIDRv4,
		CIDRv6: net.IPNet{},
		MTU:    1600,
	}, appContext); err != nil {
		return err
	} else {
		punDev = dev
	}
	return nil
}

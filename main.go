package main

import (
	"Pier/common"
	"Pier/core/service"
	"Pier/core/transport"
	"Pier/core/tunnel"
	"Pier/internal"
	"context"
	"log"
)

var appContext = context.Background()
var transportImpl = new(transport.TcpTransport)
var tunDev = new(tunnel.Device)

func main() {
	var devName string
	if config.Extra.ObfName {
		devName = common.ObfuscateText(config.TunName)
	} else {
		devName = config.TunName
	}
	if err := tunDev.Init(devName, config.CIDRv4, appContext); err != nil {
		log.Fatal("[Initial]Fail to init tunDevice:", err)
		return
	}
	log.Println("[Initial]Init tunDevice success!")
	defer tunDev.Destroy()
	if config.HttpService.Enable {
		go func(dev *tunnel.Device) {
			if err := service.StartHttpServer(config.HttpService.Port, config.HttpService.ServiceKey, dev); err != nil {
				log.Printf("[HttpService]Listen %d fail: %v\n", config.HttpService.Port, err)
			}
		}(tunDev)
	}
	if config.ServerMode {
		if err := internal.PierServer(config.ServiceAddr, tunDev, appContext, transportImpl); err != nil {
			log.Fatal("[Server]Fail to start Service.", err)
			return
		}
	} else {
		if err := internal.PierClient(config.ServiceAddr, config.CIDRv4, tunDev, appContext, transportImpl); err != nil {
			log.Fatal("[Client]Fail to start Client.", err)
			return
		}
	}
	log.Println("Unknown State.Waiting signal...")
	common.WaitingSignal()
}

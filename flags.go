package main

import (
	"Pier/common"
	"Pier/core"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var config core.Config

const PierDir = "/etc/pier"
const PierLogDir = "/etc/pier/log"

func init() {
	flag.Usage = func() {
		fmt.Println("Pier: A network tool that makes two devices neighbors.")
		fmt.Println("Usage of Pier:")
		flag.PrintDefaults()
	}
	confPath := flag.String("c", filepath.Join(PierDir, "config.toml"), "config file path(toml).")
	demo := flag.Bool("demo", false, "Give a demo config file.")
	logToFile := flag.Bool("l", false, "Redirect logs to /etc/pier/log/hh-mm.log.")
	flag.Parse()

	if *demo {
		if err := common.Save("./demo.toml", core.Config{
			CIDRv4:      "10.10.10.1/24",
			TunName:     "PierTun",
			ServiceAddr: ":38324",
			ServerMode:  true,
			HttpService: core.HttpService{
				Enable:     true,
				Port:       14122,
				ServiceKey: "A-kua",
			},
			Extra: core.Extra{
				ObfName: true,
			},
		}); err != nil {
			log.Fatalln("Internal error ", err)
			return
		}
		log.Println("Demo config generated! (./demo.toml)")
		os.Exit(1)
	}
	err := os.MkdirAll(PierLogDir, 0755)
	if err != nil {
		log.Fatalf("error creating pier directory: %v", err)
	}
	if err := common.Load(*confPath, &config); err != nil {
		log.Fatalln(err)
	}
	if *logToFile {
		now := time.Now()
		logFileName := now.Format("15-04") + ".log"
		logFilePath := filepath.Join(PierLogDir, logFileName)

		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log file: %v", err)
		}

		log.SetOutput(file)
	}
}

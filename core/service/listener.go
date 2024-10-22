package service

import (
	"Pier/core/tunnel"
	"fmt"
	"github.com/inhies/go-bytesize"
	"io"
	"net/http"
	"runtime"
	"strings"
)

func checkPermission(w http.ResponseWriter, req *http.Request, serviceKey string) bool {
	key := req.Header.Get("key")
	if strings.EqualFold(key, serviceKey) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("No permission"))
		return false
	}
	return true
}

func StartHttpServer(port int, serviceKey string, device *tunnel.Device) error {
	http.HandleFunc("/ip", func(w http.ResponseWriter, req *http.Request) {
		ip := req.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = strings.Split(req.RemoteAddr, ":")[0]
		}
		resp := fmt.Sprintf("%v", ip)
		_, _ = io.WriteString(w, resp)
		runtime.Gosched()
	})

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		if !checkPermission(w, r, serviceKey) {
			return
		}
		read, write := device.GetReadWriteBytes()
		resp := fmt.Sprintf("read %v write %v", bytesize.New(float64(read)).String(), bytesize.New(float64(write)).String())
		_, _ = io.WriteString(w, resp)
		runtime.Gosched()
	})
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		return err
	}
	return nil
}

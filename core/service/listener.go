package service

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
)

func StartHttpServer(port int) error {
	http.HandleFunc("/ip", func(w http.ResponseWriter, req *http.Request) {
		ip := req.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = strings.Split(req.RemoteAddr, ":")[0]
		}
		resp := fmt.Sprintf("%v", ip)
		_, _ = io.WriteString(w, resp)
		runtime.Gosched()
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		return err
	}
	return nil
}

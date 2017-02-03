package httputil

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mkideal/pkg/option"
)

func IP(req *http.Request) string {
	ips := getproxy(req)
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		if len(rip) > 0 {
			return rip[0]
		}
	}
	ip := strings.Split(req.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func getproxy(req *http.Request) []string {
	if ips := req.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func JSONResponse(w http.ResponseWriter, status int, value interface{}, debug ...bool) error {
	var (
		b   []byte
		err error
	)
	if option.Bool(false, debug...) {
		b, err = json.MarshalIndent(value, "", "  ")
	} else {
		b, err = json.Marshal(value)
	}
	if err != nil {
		return err
	}
	return BlobResponse(w, status, "application/json;charset=utf-8", b)
}

func BlobResponse(w http.ResponseWriter, status int, contentType string, b []byte) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", contentType)
	_, err := w.Write(b)
	return err
}

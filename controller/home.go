package controller

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

func Index(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Println("Accessing Index Function")
	fmt.Fprint(w, "Dashboard ejlog-server")
	log.Printf("200")
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-REAL-IP")
	ipa, port, err := net.SplitHostPort(r.RemoteAddr)

	InfoLogger.Printf("IP REAL : %s OR %s OR %s OR %s", ip, ipa, port, err)
	fmt.Printf("IP REAL : %s OR %s OR %s OR %s", ip, ipa, port, err)

	ips := r.Header.Get("X-FORWARDED-FOR")
	fmt.Printf("IPS REAL : %s", ips)
}

func CheckIpCache(w http.ResponseWriter, r *http.Request) {
	c := cache.New(5*time.Minute, 10*time.Minute)

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ErrorLogger.Printf("Error Get Ip %s ", ip)
	}

	if ip == "::1" {
		ip = "127.0.0.1"
	}

	kanwil, found := c.Get(ip)
	if found {
		InfoLogger.Printf("Kanwil for ip %s : %s", kanwil, ip)
	}

	InfoLogger.Printf("Not Found Kanwil for ip %s ", ip)
}
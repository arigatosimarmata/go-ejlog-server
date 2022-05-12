package controller

import (
	"ejol/ejlog-server/models"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

func Index(w http.ResponseWriter, r *http.Request) {
	models.InfoLogger.Println("Accessing Index Function")
	fmt.Fprint(w, "Dashboard ejlog-server")
	log.Printf("200")
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-REAL-IP")
	ipa, port, err := net.SplitHostPort(r.RemoteAddr)

	models.InfoLogger.Printf("IP REAL : %s OR %s OR %s OR %s", ip, ipa, port, err)
	fmt.Printf("IP REAL : %s OR %s OR %s OR %s", ip, ipa, port, err)

	ips := r.Header.Get("X-FORWARDED-FOR")
	fmt.Printf("IPS REAL : %s", ips)
}

func CheckIpCache(w http.ResponseWriter, r *http.Request) {
	c := cache.New(5*time.Minute, 10*time.Minute)

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		models.ErrorLogger.Printf("Error Get Ip %s ", ip)
	}

	if ip == "::1" {
		ip = "127.0.0.1"
	}

	kanwil, found := c.Get(ip)
	if found {
		models.InfoLogger.Printf("Kanwil for ip %s : %s", kanwil, ip)
	}

	models.InfoLogger.Printf("Not Found Kanwil for ip %s ", ip)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

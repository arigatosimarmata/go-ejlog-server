package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Router *http.ServeMux
}

func (server *Server) Run(addr string) {
	server.Router = http.NewServeMux()
	server.initializeRoutes()

	fmt.Println("Start the development server at http://127.0.0.1" + addr)
	InfoLogger.Println("Start the development server at http://127.0.0.1" + addr)

	s := &http.Server{
		Addr:           addr,
		Handler:        server.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 0,
	}
	// log.Fatal(http.ListenAndServe(addr, server.Router))
	log.Fatal(s.ListenAndServe())
}

func (s *Server) initializeRoutes() {
	s.Router.HandleFunc("/", Index)
	s.Router.HandleFunc("/homepage", Homepage)
	s.Router.HandleFunc("/health-check", HealthCheckHandler)

	s.Router.HandleFunc("/check-ip-cache", CheckIpCache)

	//VERSI 1
	s.Router.HandleFunc("/v1ejlog-server/multiline-wincor", MultilineWincor)

	//VERSI 2 WITH KAFKA
	// s.Router.HandleFunc("/v2ejlog-server/kafka-multiline-wincor", KafkaMultilineWincor)

	//VERSI 3 SEPARATE FILE
	// s.Router.HandleFunc("/v3ejlog-server/multiline-wincor", V3MultilineWincor)
	s.Router.HandleFunc("/v3ejlog-server/multiline-wincor", V3MultilineWincorAppendFile)
	// s.Router.HandleFunc("/v3ejlog-server/multiline-wincor-1", V3MultilineWincor_1)
	s.Router.HandleFunc("/v3ejlog-server/multiline-wincor-1", V3MultilineWincor_1AppendHeaderIp)
	s.Router.HandleFunc("/v3ejlog-server/multiline-wincor-elastic", V3MultilineWincorElastic)
	s.Router.HandleFunc("/v3ejlog-server/multiline-wincor-elastic2", V3MultilineWincorElastic2)
	s.Router.HandleFunc("/v3ejlog-server/multiline-wincor-elastic2-write", V3MultilineWincorElastic2WriteOnly)

}

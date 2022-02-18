package main

import (
	"bufio"
	"ejol/ejlog-server/controller"
	"ejol/ejlog-server/job"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var server = controller.Server{}

func main_withjobexportATMMAPPING() {
	err := godotenv.Load(".env")
	if err != nil {
		controller.ErrorLogger.Fatal("Error load file env : ", err)
	}
	go job.JobCacheAtmMappings()
	go job.JobExportCountAtm()
	controller.ConsumeFileEjol()
}

func main_withJob() {
	go job.JobCacheAtmMappings()
	go controller.ConsumeFileEjol()

	server.Run(":9998")

}

func servers() {
	listener, err := net.Listen("tcp", ":9988")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("server", err)
		os.Exit(1)
	}
	data := make([]byte, 1)
	if _, err := conn.Read(data); err != nil {
		log.Fatal("server", err)
	}

	conn.Close()
}

func client() {
	conn, err := net.Dial("tcp", "localhost:9988")
	if err != nil {
		log.Fatal("client", err)
	}

	// write to make the connection closed on the server side
	if _, err := conn.Write([]byte("a")); err != nil {
		log.Printf("client: %v", err)
	}

	time.Sleep(1 * time.Second)

	// write to generate an RST packet
	if _, err := conn.Write([]byte("b")); err != nil {
		log.Printf("client: %v", err)
	}

	time.Sleep(1 * time.Second)

	// write to generate the broken pipe error
	if _, err := conn.Write([]byte("c")); err != nil {
		log.Printf("client: %v", err)
		if errors.Is(err, syscall.EPIPE) {
			log.Print("This is broken pipe error")
		}
	}
}

func main_checkpipe() {
	go servers()

	time.Sleep(3 * time.Second) // wait for server to run

	client()
}

//SO FAR WORKING WELL
func main_DEFINITION_OF_DONE() {
	err := godotenv.Load(".env")
	if err != nil {
		controller.ErrorLogger.Fatal("Error load file env : ", err)
	}
	go job.JobCacheAtmMappings()
	time.Sleep(2 * time.Second)
	server.Run(":9898")
}

func testingScanner() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case "atm-mapping-cache":
			fmt.Println("Execute scheduler.")
		case "exit":
			os.Exit(0)
		default:
			fmt.Println(line) // Println will add back the final '\n'
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

//CHECKING WITH ELASTIC SEARCH
func main() {
	// controller.ExampleElasticSearch()
	err := godotenv.Load(".env")
	if err != nil {
		controller.ErrorLogger.Fatal("Error load file env : ", err)
	}
	go job.JobCacheAtmMappings()
	time.Sleep(2 * time.Second)

	go testingScanner()

	server.Run(":3000")
	// controller.Testingduludeh2()
}

func mainasdfasdf() {
	datedetail := time.Now()
	d := time.Now().Format("150405.0000000")
	fmt.Println(datedetail)
	fmt.Println(d)
}

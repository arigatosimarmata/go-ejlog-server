package main

import (
	"ejol/ejlog-server/controller"
	"ejol/ejlog-server/job"
	"time"

	"github.com/joho/godotenv"
)

var server = controller.Server{}

//CHECKING WITH ELASTIC SEARCH
func main() {
	// controller.ExampleElasticSearch()
	err := godotenv.Load(".env")
	if err != nil {
		controller.ErrorLogger.Fatal("Error load file env : ", err)
	}
	go job.JobCacheAtmMappings()
	// go job.JobExportCountAtm()
	time.Sleep(2 * time.Second)

	server.Run(":3000")
}

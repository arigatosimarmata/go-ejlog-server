package controller

import (
	"context"
	"ejol/ejlog-server/models"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

func WincorParseProcess() error {
	date := time.Now().Format("20060102")
	dirPath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + date + "/"
	sectionDirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		models.ErrorLogger.Printf("Error Get Dir : %s", err)
		return err
	}

	for _, folder := range sectionDirs {
		fmt.Println(folder.Name())
		readFile, err := ioutil.ReadDir(dirPath + folder.Name() + "/")
		if err != nil {
			models.ErrorLogger.Printf("Error read Directory : %s", err)
			return err
		}

		for _, f := range readFile {
			fmt.Println(f.Name())

			filePath := dirPath + folder.Name() + "/"
			ej_content, err := os.ReadFile(filePath + f.Name())
			if err != nil {
				models.ErrorLogger.Printf("Error read file : %s", err)
				return err
			}

			filename := strings.Split(f.Name(), "_")
			ip_address := filename[0]
			kanwil := filename[1]

			fmt.Printf("Split ip_address : %s - kanwil : %s ", ip_address, kanwil)
			err = WincorProcessEjlogElastic(string(ej_content), ip_address, kanwil)
			if err != nil {
				models.ErrorLogger.Printf("Error processEjlogElastic : %s", err)
				return err
			}
			models.InfoLogger.Printf("Sukses menyimpan RequestEjlog dalam file %s", filePath+f.Name())

			err = os.Rename(filePath+f.Name(), filePath+f.Name())
			if err != nil {
				models.ErrorLogger.Printf("Error pada merubah file : %s", err)
				return err
			}

			models.InfoLogger.Printf("Sukses rename File %s", filePath+f.Name())
		}

	}

	return nil
}

func WincorProcessEjlogElastic(ejcontent, ip_address, kanwil string) error {
	keywordMapWincor := models.KeywordEjol
	date := time.Now()
	requestBody := ejcontent
	ejol_map := strings.Split(string(requestBody), "\n")

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		models.ErrorLogger.Fatal("Error creating the client.")
		return err
	}

	//insert data
	for _, ejlog := range ejol_map {
		space2 := regexp.MustCompile(`\s+`)
		cleanejlog := space2.ReplaceAllString(ejlog, " ")

		for key, value := range (*keywordMapWincor)["WINCOR_KEYWORD"] {
			if strings.Contains(cleanejlog, value) {
				req := esapi.IndexRequest{
					Index: "ejmonitor-" + date.Format("20060102"),
					Body: strings.NewReader(`
					{
						"ejlog" : "` + cleanejlog + `",
						"ip_address" :"` + ip_address + `",
						"parameter":"` + key + `",
						"tanggal":"` + time.Now().Format("20060102") + `",
					}
					`),
					Refresh: "true",
				}

				res, err := req.Do(context.Background(), es)
				if err != nil {
					log.Fatalf("Error getting response: %s", err)
					return err
				}
				defer res.Body.Close()

				if res.IsError() {
					log.Printf("[%s] Error indexing document ", res.Status())
					return err
				}
			}
		}
	}

	return nil
}

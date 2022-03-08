package controller

import (
	"context"
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

func HyosungParseProcessByAtmMapping() error {
	date := time.Now().Format("20060102")
	dirPath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + date + "/"
	sectionDirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		ErrorLogger.Printf("Error Get Dir : %s", err)
		return err
	}

	for _, folder := range sectionDirs {
		fmt.Println(folder.Name())
		readFile, err := ioutil.ReadDir(dirPath + folder.Name() + "/")
		if err != nil {
			ErrorLogger.Printf("Error read Directory : %s", err)
			return err
		}

		for _, f := range readFile {
			fmt.Println(f.Name())

			filePath := dirPath + folder.Name() + "/"
			ej_content, err := os.ReadFile(filePath + f.Name())
			if err != nil {
				ErrorLogger.Printf("Error read file : %s", err)
				return err
			}

			filename := strings.Split(f.Name(), "_")
			ip_address := filename[0]
			kanwil := filename[1]

			fmt.Printf("Split ip_address : %s - kanwil : %s ", ip_address, kanwil)
			err = processEjlogElastic(string(ej_content), ip_address, kanwil)
			if err != nil {
				ErrorLogger.Printf("Error processEjlogElastic : %s", err)
				return err
			}
			InfoLogger.Printf("Sukses menyimpan RequestEjlog dalam file %s", filePath+f.Name())

			err = os.Rename(filePath+f.Name(), filePath+f.Name())
			if err != nil {
				ErrorLogger.Printf("Error pada merubah file : %s", err)
				return err
			}

			InfoLogger.Printf("Sukses rename File %s", filePath+f.Name())
		}

	}

	return nil
}

func HyosungParseProcess() error {
	date := time.Now().Format("20060102")
	dirPath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + date + "/"
	sectionDirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		ErrorLogger.Printf("Error Get Dir : %s", err)
		return err
	}

	for _, folder := range sectionDirs {
		fmt.Println(folder.Name())
		readFile, err := ioutil.ReadDir(dirPath + folder.Name() + "/")
		if err != nil {
			ErrorLogger.Printf("Error read Directory : %s", err)
			return err
		}

		for _, f := range readFile {
			fmt.Println(f.Name())

			filePath := dirPath + folder.Name() + "/"
			ej_content, err := os.ReadFile(filePath + f.Name())
			if err != nil {
				ErrorLogger.Printf("Error read file : %s", err)
				return err
			}

			filename := strings.Split(f.Name(), "_")
			ip_address := filename[0]
			kanwil := filename[1]

			fmt.Printf("Split ip_address : %s - kanwil : %s ", ip_address, kanwil)
			err = processEjlogElastic(string(ej_content), ip_address, kanwil)
			if err != nil {
				ErrorLogger.Printf("Error processEjlogElastic : %s", err)
				return err
			}
			InfoLogger.Printf("Sukses menyimpan RequestEjlog dalam file %s", filePath+f.Name())

			err = os.Rename(filePath+f.Name(), filePath+f.Name())
			if err != nil {
				ErrorLogger.Printf("Error pada merubah file : %s", err)
				return err
			}

			InfoLogger.Printf("Sukses rename File %s", filePath+f.Name())
		}

	}

	return nil
}

func processEjlogElastic(ejcontent, ip_address, kanwil string) error {
	keywordMapHyosung := KeywordEjol

	date := time.Now()
	requestBody := ejcontent
	ejol_map := strings.Split(string(requestBody), "\n")

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		ErrorLogger.Fatal("Error creating the client.")
		return err
	}

	//insert data
	for _, ejlog := range ejol_map {
		space2 := regexp.MustCompile(`\s+`)
		cleanejlog := space2.ReplaceAllString(ejlog, " ")

		for key, value := range (*keywordMapHyosung)["HYOSUNG_KEYWORD"] {
			if strings.Contains(cleanejlog, value) {
				req := esapi.IndexRequest{
					Index: "ejmonitor-" + date.Format("20060102"),
					Body: strings.NewReader(`
					{"ejlog" : "` + cleanejlog + `",
					"ip_address" :"` + ip_address + `",
					"keyword":"` + key + `",
					"date":"` + time.Now().Format("2006-01-02 15:04:05") + `"}
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

// func processEjlogElastic2(ejcontent, ip_address, kanwil string) error {
// 	date := time.Now()
// 	requestBody := ejcontent
// 	ejol_map := strings.Split(string(requestBody), "\n")

// 	es, err := elasticsearch.NewDefaultClient()
// 	if err != nil {
// 		ErrorLogger.Fatal("Error creating the client.")
// 		return err
// 	}

// 	//insert data
// 	for _, ejlog := range ejol_map {
// 		space2 := regexp.MustCompile(`\s+`)
// 		cleanejlog := space2.ReplaceAllString(ejlog, " ")

// 		for _, value := range keywordMapHyosung {
// 			if strings.Contains(cleanejlog, value) {
// 				req := esapi.IndexRequest{
// 					Index: "ejmonitor-" + date.Format("20060102"),
// 					Body: strings.NewReader(`
// 					{"ejlog" : "` + cleanejlog + `",
// 					"ip_address" :"` + ip_address + `"}
// 					`),
// 					Refresh: "true",
// 				}

// 				res, err := req.Do(context.Background(), es)
// 				if err != nil {
// 					log.Fatalf("Error getting response: %s", err)
// 					return err
// 				}
// 				defer res.Body.Close()

// 				if res.IsError() {
// 					log.Printf("[%s] Error indexing document ", res.Status())
// 					return err
// 				}
// 			}
// 		}
// 		// if strings.Contains(cleanejlog, HyosungParser()("DISPENSE_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("PRINT_CASH_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("ERROR_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("CARD_JAMMED_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("PRINT_CASH_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("CARD_RETAIN_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("COMMUNICATION_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("OOS_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("OOS_KEYWORD2")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("INSERVICE_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("INSERVICE_KEYWORD2")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("RECEIPT_PRINTER_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("RECEIPT_PAPER_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// } else if strings.Contains(cleanejlog, HyosungParser()("TRANSACTION_KEYWORD")) {
// 		// 	req := esapi.IndexRequest{
// 		// 		Index: "ejmonitor-" + date.Format("20060102"),
// 		// 		Body: strings.NewReader(`
// 		// 		{"ejlog" : "` + cleanejlog + `",
// 		// 		"ip_address" :"` + ip_address + `"}
// 		// 		`),
// 		// 		Refresh: "true",
// 		// 	}

// 		// 	res, err := req.Do(context.Background(), es)
// 		// 	if err != nil {
// 		// 		log.Fatalf("Error getting response: %s", err)
// 		// 		return err
// 		// 	}
// 		// 	defer res.Body.Close()

// 		// 	if res.IsError() {
// 		// 		log.Printf("[%s] Error indexing document ", res.Status())
// 		// 		return err
// 		// 	}
// 		// }

// 	}

// 	return nil
// }

func HyosungParseErrorKeyword() {

}

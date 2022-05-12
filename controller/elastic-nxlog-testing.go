package controller

import (
	"context"
	"ejol/ejlog-server/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

func V3MultilineWincorElastic(w http.ResponseWriter, r *http.Request) {
	date := time.Now()
	ip_address, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		models.ErrorLogger.Printf("RC : %d - Error %s", http.StatusNotAcceptable, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if ip_address == "::1" {
		ip_address = "127.0.0.1"
	}

	getKanwil, found := Cac.Get(ip_address)
	if !found {
		models.ErrorLogger.Printf("RC : %d - Ip Not Found ", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	kanwil := getKanwil.(string)
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		models.ErrorLogger.Fatal("Error creating the client.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	requestBody, _ := ioutil.ReadAll(r.Body)
	ejol_map := strings.Split(string(requestBody), "\n")
	for i, ejlog := range ejol_map {
		space2 := regexp.MustCompile(`\s+`)
		cleanejlog := space2.ReplaceAllString(ejlog, " ")
		req := esapi.IndexRequest{
			Index: "ejlog_" + kanwil + "_" + date.Format("20060102"),
			// DocumentID: strconv.Itoa(i + 1),
			Body: strings.NewReader(`
			{"ejlog" : "` + cleanejlog + `",
			"ip_address" :"` + ip_address + `"}
			`),
			Refresh: "true",
		}

		log.Printf("Data ke %d - contoh request : %s", i, req.Body)
		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("[%s] Error indexing document ID=%d", res.Status(), i)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			// Deserialize the response into a map.
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Printf("Error parsing the response body: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				// Print the response status and indexed document version.
				log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

func V3MultilineWincorElastic(w http.ResponseWriter, r *http.Request) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal("Error creating the client")
	}

	requestBody, _ := ioutil.ReadAll(r.Body)
	ejol_map := strings.Split(string(requestBody), "\n")
	ip_address := "127.0.0.1"
	for i, ejlog := range ejol_map {
		// log.Printf("Data ke %d --- Ejlog - %s", i+1, ejlog)

		space2 := regexp.MustCompile(`\s+`)
		cleanejlog := space2.ReplaceAllString(ejlog, " ")
		req := esapi.IndexRequest{
			Index: "ejlog_20_" + time.Now().Format("20060102"),
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

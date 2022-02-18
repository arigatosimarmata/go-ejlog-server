package controller

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

func Testingduludeh2(w http.ResponseWriter, r *http.Request) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal("Error creating the client")
	}

	date := time.Now().Format("2006-01-02 15:04:05.000")
	res, err := es.Index(
		"test-myindex-0001",
		strings.NewReader(`{
		"ip_address": "ej_127_0_0_1",
		"ejlog":[
			{
				"ejlog": "1. programming list"
			},
			{
				"ejlog": "2. cool_list"
			},
			{"ejlog": "3. cool_stuff_list"},
			{"ejlog": "4. askdfajsdlf"},
			{"ejlog": "5. programming_list"},
			{"ejlog": "6. cool_list"},
			{"ejlog": "7. cool_stuff_list"},
			{"ejlog": "8. askdfajsdlf"},
			{"ejlog": "9. programming_list"},
			{"ejlog": "10. cool_list"},
			{"ejlog": "11. cool_stuff_list"},
			{"ejlog": "12. askdfajsdlf"}
		],
		"date" : "`+date+`"
	}`),
		es.Index.WithDocumentID("20"),
		es.Index.WithPretty(),
	)
	// fmt.Println(res, err)

	// res, err := es.Search(
	// 	es.Search.WithContext(context.Background()),
	// 	es.Search.WithIndex("test-myindex-0001"),
	// 	es.Search.WithBody(strings.NewReader(`{"query" : { "match" : { "ip_address" : "ej_127_0_0_1" } }}`)),
	// 	es.Search.WithTrackTotalHits(true),
	// 	es.Search.WithPretty(),
	// )
	// if err != nil {
	// 	log.Fatalf("ERROR: %s", err)
	// }
	// defer res.Body.Close()

	// if res.IsError() {
	// 	var e map[string]interface{}
	// 	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
	// 		log.Fatalf("error parsing the response body: %s", err)
	// 	} else {
	// 		// Print the response status and error information.
	// 		log.Fatalf("[%s] %s: %s",
	// 			res.Status(),
	// 			e["error"].(map[string]interface{})["type"],
	// 			e["error"].(map[string]interface{})["reason"],
	// 		)
	// 	}
	// }

	// // var r map[string]interface{}
	// // if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
	// // 	log.Fatalf("Error parsing the response body: %s", err)
	// // }
	// // // Print the response status, number of results, and request duration.
	// // log.Printf(
	// // 	"[%s] %d hits; took: %dms",
	// // 	res.Status(),
	// // 	int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
	// // 	int(r["took"].(float64)),
	// // )
	// // // Print the ID and document source for each hit.
	// // for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
	// // 	log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	// // }

	fmt.Println(res, err)
}

func Testingduludeh(w http.ResponseWriter, r *http.Request) {
	// cfg := elasticsearch.Config{
	// 	Addresses: []string{
	// 		"http://127.0.0.1:9200",
	// 	},
	// 	Transport: &http.Transport{
	// 		MaxIdleConnsPerHost:   9,
	// 		ResponseHeaderTimeout: time.Second,
	// 		DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
	// 		TLSClientConfig: &tls.Config{
	// 			MaxVersion:         tls.VersionTLS11,
	// 			InsecureSkipVerify: true,
	// 		},
	// 	},
	// }

	// es, err := elasticsearch.NewClient(cfg)
	// if err != nil {
	// 	log.Fatal("Error creating the client")
	// 	// w.WriteHeader(http.StatusInternalServerError)
	// 	// return
	// }
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal("Error creating the client")
	}

	requestBody, _ := ioutil.ReadAll(r.Body)
	// ejol_map := strings.Split(string(requestBody), "/\\r\\n|\\r|\\n|\x0D\x0A|\n/")
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

		// Perform the request with the client.

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

func V3MultilineWincorElastic(w http.ResponseWriter, r *http.Request) {
	// var b bytes.Buffer
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal("Error creating the client")
	}

	requestBody, _ := ioutil.ReadAll(r.Body)
	// ejol_map := strings.Split(string(requestBody), "\n")

	// for _, ejlog := range ejol_map {
	// req := esapi.IndexRequest{
	// 	Index:      "ejlog_20_20220210",
	// 	DocumentID: strconv.Itoa(i + 1),
	// 	Body:       strings.NewReader(`{"ejlog" : "` + ejlog + `"}`),
	// 	Refresh:    "true",
	// }

	req := esapi.BulkRequest{
		Index:   "ejlog_20_20220210",
		Body:    strings.NewReader(`{"ejlog" : "` + string(requestBody) + `"}`),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response : %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ", res.Status())
	} else {
		var rsp map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&rsp); err != nil {
			log.Printf("Error parsing the response body : %s", err)
		} else {
			log.Printf("[%s] %s; ", res.Status(), rsp["result"])
		}
	}

	// b.WriteString(`{"ejlog" : "`)
	// b.WriteString(ejlog)
	// b.WriteString(`"}`)
	// req, _ := es.Index("ejlog_20_20220210", &b)
	// fmt.Println("-")

	// }
}

func MultilineWincorElastic_ConfigClient() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTIC_HOST") + ":" + os.Getenv("ELASTIC_PORT"),
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MaxVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal("Error creating the client")
	}

	fmt.Println(es)

}

func ExampleElasticSearch() {
	log.SetFlags(0)

	var (
		r  map[string]interface{}
		wg sync.WaitGroup
	)

	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 1. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print version number.
	log.Printf("~~~~~~~> Elasticsearch %s", r["version"].(map[string]interface{})["number"])

	// 2. Index documents concurrently
	//
	for i, title := range []string{"23:38:40 PIN ENTERED",
		"23:38:42 ***** Tran Request State *****",
		"23:38:42 TRANSACTION REQUEST AF      ",
		"23:38:43 Transaction reply next 523 function 5157",
		"23:38:43 TVR: 8000040000, TSI: 7000",
		"TANGGAL: 08/12/21   WAKTU : 23:57:56",
		"ATM ID : 790121     NO.REF: 08687",
		"CARD   : 522184XXXXXX2943    ",
		"REK    : 4067406701010928507",
		"KODETX : 99   CRM ",
		"AMOUNT : RP 0",
		"SALDO  : RP 0",
		"RESPON : W0   INQUIRY",
		"23:38:43 ***** Tran Request State *****",
		"23:38:44 ***** Accept Cash State *****",
		"23:38:44 $MOD$ 210216_1011 BriCRMTa.DLL",
		"23:38:44 ----- Clear Cash-In Device ----",
		"23:38:44 ----- Insert Notes ----",
		"23:38:53   Notes INSERTED",
		"23:38:53 ----- Identify Notes ----",
		"23:39:00 INVALID NOTES RECOGNIZED",
		"23:39:00 TYPE1TYPE2",
		"23:39:00   Notes IDENTIFIED : (IDR,50000,6)",
		"23:39:01   Notes IDENTIFIED : (IDR,100000,1)",
		"23:39:01   Notes REFUSED : 1",
		"23:39:01 ----- Remove Notes ----",
		"23:39:10   Refused Notes in Output Tray TAKEN ",
		"23:39:12 ----- Confirm Notes ----",
		"23:39:14   User Pressed DEPOSIT Button",
		"23:39:14 ***** Accept Cash State *****",
		"23:39:16 EMV AID A0000006021010 / 522184******2943 STARTED"} {
		wg.Add(1)

		go func(i int, title string) {
			fmt.Printf("data %d - content : %s", i, title)
			defer wg.Done()

			// Set up the request object directly.
			req := esapi.IndexRequest{
				Index:      "test",
				DocumentID: strconv.Itoa(i + 1),
				Body:       strings.NewReader(`{"ejlog" : "` + title + `"}`),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(i, title)
	}
	wg.Wait()

	log.Println(strings.Repeat("-", 37))

	// 3. Search for the indexed documents
	//
	// Use the helper methods of the client.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithBody(strings.NewReader(`{"query" : { "match" : { "title" : "test" } }}`)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))
}

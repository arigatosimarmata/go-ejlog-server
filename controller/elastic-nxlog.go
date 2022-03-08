package controller

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch"
)

type EjModel struct {
	Ejlog string `json:"ejlog"`
}

func V3MultilineWincorElastic2(w http.ResponseWriter, r *http.Request) {
	date := time.Now()
	ejlog_arrays := []*EjModel{}
	ip_address, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ErrorLogger.Printf("RC : %d - Error %s", http.StatusNotAcceptable, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if ip_address == "::1" {
		ip_address = "127.0.0.1"
	}

	getKanwil, found := Cac.Get(ip_address)
	if !found {
		ErrorLogger.Printf("RC : %d - Ip Not Found ", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	kanwil := getKanwil.(string)
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		ErrorLogger.Fatal("Error creating the client.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	requestBody, _ := ioutil.ReadAll(r.Body)
	ejol_map := strings.Split(string(requestBody), "\n")

	for _, ejlog := range ejol_map {
		space2 := regexp.MustCompile(`\s+`)
		cleanejlog := space2.ReplaceAllString(ejlog, " ")
		ejlog_array := &EjModel{Ejlog: cleanejlog}
		ejlog_arrays = append(ejlog_arrays, ejlog_array)
	}

	ejmap_string := new(strings.Builder)
	json.NewEncoder(ejmap_string).Encode(ejlog_arrays)

	res, err := es.Index(
		"ejlog_"+kanwil+"_"+date.Format("20060102"),
		strings.NewReader(`{
			"ip_address" : "`+ip_address+`",
			"ejlog" : `+ejmap_string.String()+`,
			"date" : "`+date.Format("2006-01-02 15:04:05.000")+`"
			}`),
		// es.Index.WithDocumentID("1"),
		es.Index.WithPretty(),
		es.Index.WithRefresh("true"),
	)

	if err != nil {
		ErrorLogger.Printf("Error getting response : %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// fmt.Println(res, err)
	InfoLogger.Println(res.Status(), err)
	w.WriteHeader(http.StatusOK)
}

func V3MultilineWincorElastic2WriteOnly(w http.ResponseWriter, r *http.Request) {
	ejlog_arrays := []*EjModel{}
	requestBody, _ := ioutil.ReadAll(r.Body)
	ejol_map := strings.Split(string(requestBody), "\n")

	for _, ejlog := range ejol_map {
		space2 := regexp.MustCompile(`\s+`)
		cleanejlog := space2.ReplaceAllString(ejlog, " ")
		ejlog_array := &EjModel{Ejlog: cleanejlog}
		ejlog_arrays = append(ejlog_arrays, ejlog_array)
	}

	ejmap_string := new(strings.Builder)
	json.NewEncoder(ejmap_string).Encode(ejlog_arrays)
	InfoLogger.Println(ejmap_string.String())
	w.WriteHeader(http.StatusOK)
}

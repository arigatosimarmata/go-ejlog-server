//NXLOG VERSI 3 MENYIMPAN DALAM BENTUK FILE
package controller

import (
	"bufio"
	"ejol/ejlog-server/config"
	"ejol/ejlog-server/models"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/miguelmota/go-filecache"
	"github.com/patrickmn/go-cache"
)

var Cac *cache.Cache

func V3MultilineWincor_1(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	start := time.Now()

	todayDate := start.Format("20060102")
	ip_address, _, error := net.SplitHostPort(r.RemoteAddr)

	if error != nil {
		models.ErrorLogger.Printf("RC : %d - Error : %s", http.StatusNotAcceptable, error)
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
	tblname := strings.ReplaceAll(ip_address, ".", "_")

	// storePath := "appendrow/" + todayDate + "/" + tblname
	storePath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate + "/" + tblname

	if _, errPath := os.Stat(storePath); os.IsNotExist(errPath) {
		err := os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			models.ErrorLogger.Printf("Error : %s", err)
			return
		}
	}

	if r.Header.Get("Content-Type") == "" {
		log.Printf("Error Content Type empty, %v", r.Header.Get("Content-Type"))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not read received POST payload: %v", err)))
		return
	}
	content := string(rb)
	headerName := ip_address + "_" + strings.ReplaceAll(start.Format("150405.0000000"), ".", "")
	filename := storePath + "/" + headerName
	err = ioutil.WriteFile(filename, []byte(content), 0666)
	if err != nil {
		models.ErrorLogger.Printf("RC : %d - Error WriteFile : %s", http.StatusNotAcceptable, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	ejdata, err := processRequestEjlog2(content, ip_address, kanwil)
	if err != nil {
		models.ErrorLogger.Printf("RC : %d - Error processRequestEjlog2 : %s, ejdata : %s ", http.StatusNotAcceptable, err, ejdata)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	models.InfoLogger.Printf("RC : %d - Save successfully.", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

func V3MultilineWincor_1AppendHeaderIp(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	start := time.Now()

	todayDate := start.Format("20060102")
	ip_address := r.Header.Get("IP-ADDRESS")

	if ip_address == "::1" {
		ip_address = "127.0.0.1"
	}

	getKanwil, found := Cac.Get(ip_address)
	if !found {
		models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).Str("ip_address", ip_address).Msg("Not Found in Cache")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	kanwil := getKanwil.(string)

	storePath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate
	if r.Header.Get("Content-Type") == "" {
		log.Printf("Error Content Type empty, %v", r.Header.Get("Content-Type"))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not read received POST payload: %v", err)))
		return
	}
	content := string(rb)
	headerName := "ej_" + strings.ReplaceAll(ip_address, ".", "_")
	filename := storePath + "/" + headerName

	if _, errPath := os.Stat(storePath); os.IsNotExist(errPath) {
		err := os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).
				Str("ip_address", ip_address).
				Err(err)
			return
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).
			Str("ip_address", ip_address).
			Err(err)
		return
	}
	defer file.Close()
	if _, err := file.WriteString(content + "\n"); err != nil {
		models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).
			Str("ip_address", ip_address).
			Err(err)
		return
	}

	ejdata, err := processRequestEjlog2(content, ip_address, kanwil)
	if err != nil {
		models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).
			Int("rc", http.StatusNotAcceptable).
			Str("ip_address", ip_address).
			Str("content", ejdata).
			Err(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	elapsed := time.Since(start)
	models.Logger.Info().Str("time", time.Now().Format("2006-01-02T15:04:05")).
		Int("rc", http.StatusOK).
		Str("ip_address", ip_address).
		Str("content", ejdata).
		Str("process_time", elapsed.String()).
		Msg("Success.")
	w.WriteHeader(http.StatusOK)
}

func DebugV3MultilineWincor_1AppendHeaderIp(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	start := time.Now()

	todayDate := start.Format("20060102")
	ip_address := r.Header.Get("IP-ADDRESS")

	if ip_address == "::1" {
		ip_address = "127.0.0.1"
	}

	getKanwil, found := Cac.Get(ip_address)
	if !found {
		// models.ErrorLogger.Printf("RC : %d - Ip : %s Not Found ", http.StatusNotFound, ip_address)
		models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).Str("ip_address", ip_address).Msg("Not Found in Cache")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	kanwil := getKanwil.(string)

	storePath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate
	if r.Header.Get("Content-Type") == "" {
		log.Printf("Error Content Type empty, %v", r.Header.Get("Content-Type"))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not read received POST payload: %v", err)))
		return
	}
	content := string(rb)
	headerName := "ej_" + strings.ReplaceAll(ip_address, ".", "_")
	filename := storePath + "/" + headerName

	if _, errPath := os.Stat(storePath); os.IsNotExist(errPath) {
		err := os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			// models.ErrorLogger.Printf("Ip : %s Error : %s", ip_address, err)
			models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).Str("ip_address", ip_address).Err(err)
			return
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// models.ErrorLogger.Printf("Ip : %s Error : %s", ip_address, err)
		models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).Str("ip_address", ip_address).Err(err)
		return
	}
	defer file.Close()
	if _, err := file.WriteString(content + "\n"); err != nil {
		// models.ErrorLogger.Fatal(err)
		models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).Str("ip_address", ip_address).Err(err)
		return
	}

	ejdata, err := processRequestEjlog2(content, ip_address, kanwil)
	if err != nil {
		// models.ErrorLogger.Printf("RC : %d - %s Error processRequestEjlog2 : %s", http.StatusNotAcceptable, ip_address, err)
		models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).
			Int("rc", http.StatusNotAcceptable).
			Str("ip_address", ip_address).
			Str("content", ejdata).
			Err(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	// models.InfoLogger.Printf("RC : %d - %s", http.StatusOK, ip_address)
	models.Logger.Info().Str("time", time.Now().Format("2006-01-02T15:04:05")).Int("rc", http.StatusOK).Str("ip_address", ip_address).Msg("Success.")
	w.WriteHeader(http.StatusOK)
}

func V3MultilineWincorAppendFile(w http.ResponseWriter, r *http.Request) {
	ip_address := r.Header.Get("IP-ADDRESS")

	if ip_address == "::1" {
		ip_address = "127.0.0.1"
	}

	requestBody, _ := ioutil.ReadAll(r.Body)
	ejol_map := strings.Split(string(requestBody), "\n")
	namefile := ejol_map[0]

	kanwil, found := Cac.Get(ip_address)
	if !found {
		models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).Int("rc", http.StatusNotFound).Str("ip_address", ip_address).Msg("Cache Not Found.")
		err := UnmappingAgentByNxLogAppendFile(ip_address, string(requestBody))
		if err != nil {
			models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).
				Str("section", "Unmapping Agent").
				Int("rc", http.StatusInternalServerError).
				Str("ip_address", ip_address).
				Str("content", fmt.Sprint(ejol_map)).
				Err(err)
			w.WriteHeader(http.StatusOK)
			return
		}
		models.Logger.Info().Str("time", time.Now().Format("2006-01-02T15:04:05")).
			Str("section", "Agent CURL").
			Int("rc", http.StatusOK).Str("ip_address", ip_address).
			Str("content", fmt.Sprint(ejol_map)).
			Msg("Success.")
		w.WriteHeader(http.StatusOK)
		return
	} else {

		if strings.Contains(namefile, strings.ToUpper("TRANSFER_OUTPUT")) {
			err := AgentByCURL(ip_address, string(requestBody), ejol_map[0], kanwil.(string))
			if err != nil {
				models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).
					Str("section", "Agent CURL").
					Int("rc", http.StatusInternalServerError).
					Str("ip_address", ip_address).
					Str("content", fmt.Sprint(ejol_map)).
					Err(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			models.Logger.Info().Str("time", time.Now().Format("2006-01-02T15:04:05")).
				Str("section", "Agent CURL").
				Int("rc", http.StatusOK).
				Str("ip_address", ip_address).
				Str("content", fmt.Sprint(ejol_map)).
				Msg("Success.")
			return
		} else {
			err := AgentByNxLogAppendFile(ip_address, string(requestBody), kanwil.(string))
			if err != nil {
				models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).
					Str("section", "Agent NxLog").
					Int("rc", http.StatusInternalServerError).
					Str("ip_address", ip_address).
					Str("content", fmt.Sprint(ejol_map)).
					Err(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			models.Logger.Info().Str("time", time.Now().Format("2006-01-02T15:04:05")).
				Str("section", "Agent CURL").
				Int("rc", http.StatusOK).Str("ip_address", ip_address).
				Str("content", fmt.Sprint(ejol_map)).
				Msg("Success.")
			return
		}
	}

}

func DebugV3MultilineWincorAppendFile(w http.ResponseWriter, r *http.Request) {
	ip_address := r.Header.Get("IP-ADDRESS")

	if ip_address == "::1" {
		ip_address = "127.0.0.1"
	}

	kanwil, found := Cac.Get(ip_address)
	if !found {
		// models.ErrorLogger.Printf("RC : %d - Ip %s Not Found ", http.StatusNotFound, ip_address)
		models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).Int("rc", http.StatusNotFound).Str("ip_address", ip_address).Msg("Cache Not Found.")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Header.Get("Content-Type") == "" {
		log.Printf("Error Content Type empty, %v", r.Header.Get("Content-Type"))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not read received POST payload: %v", err)))
		return
	}
	ejol_map := strings.Split(string(requestBody), "\n")
	namefile := ejol_map[0]

	if strings.Contains(namefile, strings.ToUpper("TRANSFER_OUTPUT")) {
		fmt.Println("inCURl")
		err := AgentByCURL(ip_address, string(requestBody), ejol_map[0], kanwil.(string))
		if err != nil {
			// models.ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).Str("section", "Agent CURL").Int("rc", http.StatusInternalServerError).Str("ip_address", ip_address).Err(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		models.Logger.Info().Str("time", time.Now().Format("2006-01-02T15:04:05")).Str("section", "Agent CURL").Int("rc", http.StatusOK).Str("ip_address", ip_address).Msg("Success.")
		return
	} else {
		fmt.Println("inAgent")
		err := AgentByNxLogAppendFile(ip_address, string(requestBody), kanwil.(string))
		if err != nil {
			// models.ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			models.Logger.Error().Str("time", time.Now().Format("2006-01-02T15:04:05")).Str("section", "Agent NxLog").Int("rc", http.StatusInternalServerError).Str("ip_address", ip_address).Err(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		models.Logger.Info().Str("time", time.Now().Format("2006-01-02T15:04:05")).Str("section", "Agent CURL").Int("rc", http.StatusOK).Str("ip_address", ip_address).Msg("Success.")
		return
	}
}

func V3MultilineWincorSplitFile(w http.ResponseWriter, r *http.Request) {
	ip_address, _, error := net.SplitHostPort(r.RemoteAddr)
	if error != nil {
		models.ErrorLogger.Printf("RC : %d - %s Error : %s", http.StatusNotAcceptable, ip_address, error)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if ip_address == "::1" {
		ip_address = "127.0.0.1"
	}

	kanwil, found := Cac.Get(ip_address)
	if !found {
		models.ErrorLogger.Printf("RC : %d - Ip %s Not Found ", http.StatusNotFound, ip_address)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Header.Get("Content-Type") == "" {
		log.Printf("Error Content Type empty, %v", r.Header.Get("Content-Type"))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not read received POST payload: %v", err)))
		return
	}
	ejol_map := strings.Split(string(requestBody), "\n")
	namefile := ejol_map[0]

	if strings.EqualFold(namefile[0:15], strings.ToUpper("TRANSFER_OUTPUT")) {
		err := AgentByCURL(ip_address, string(requestBody), ejol_map[0], kanwil.(string))
		if err != nil {
			models.ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error %s", err)
			return
		}
	} else {
		err := AgentByNxLog(ip_address, string(requestBody), kanwil.(string))
		if err != nil {
			models.ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error %s", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	models.InfoLogger.Printf("RC : %d - %s", http.StatusOK, ip_address)
	fmt.Fprintf(w, "200")
}

func AgentByCURL(ip_address, requestBody, filename, kanwil string) error {
	start := time.Now()
	todayDate := start.Format("20060102")

	tblname := strings.ReplaceAll(ip_address, ".", "_")
	storePath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate + "/" + tblname
	// storePath := "appendrow/" + todayDate + "/" + tblname

	if _, errPath := os.Stat(storePath); os.IsNotExist(errPath) {
		err := os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			models.ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			return err
		}
	}

	content := string(requestBody)
	filesave := storePath + "/" + filename
	err := ioutil.WriteFile(filesave, []byte(content), 0666)
	if err != nil {
		models.ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
		return err
	}

	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	return nil
}

func AgentByNxLog(ip_address, requestBody, kanwil string) error {
	start := time.Now()
	todayDate := start.Format("20060102")

	tblname := strings.ReplaceAll(ip_address, ".", "_")
	storePath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate + "/" + tblname
	// storePath := "appendrow/" + todayDate + "/" + tblname

	if _, errPath := os.Stat(storePath); os.IsNotExist(errPath) {
		err := os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			models.ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			return err
		}
	}

	content := string(requestBody)
	// headerName := ip_address + "_" + strings.ReplaceAll(start.Format("150405.0000000"), ".", "")
	unixtime := time.Now().UnixNano()
	headerName := ip_address + "_" + kanwil + "_" + strconv.Itoa(int(unixtime))
	filename := storePath + "/" + headerName
	err := ioutil.WriteFile(filename, []byte(content), 0666)
	if err != nil {
		models.ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
		return err
	}

	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	return nil
}

func AgentByNxLogAppendFile(ip_address, requestBody, kanwil string) error {
	start := time.Now()
	todayDate := start.Format("20060102")

	// storePath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate + "/" + tblname
	storePath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate

	if _, errPath := os.Stat(storePath); os.IsNotExist(errPath) {
		err := os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			models.ErrorLogger.Printf("Ip : %s Error : %s", ip_address, err)
			return err
		}
	}

	content := string(requestBody)
	headerName := "ej_" + strings.ReplaceAll(ip_address, ".", "_")
	filename := storePath + "/" + headerName

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		models.ErrorLogger.Printf("Ip : %s Error : %s", ip_address, err)
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(content + "\n"); err != nil {
		models.ErrorLogger.Fatal(err)
		return err
	}

	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	return nil
}

func UnmappingAgentByNxLogAppendFile(ip_address, requestBody string) error {
	start := time.Now()
	todayDate := start.Format("20060102")
	storePath := ""

	if *models.Unmapping == "" {
		storePath = os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate + "/" + *models.Unmapping
	} else {
		storePath = os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate
	}

	if _, errPath := os.Stat(storePath); os.IsNotExist(errPath) {
		err := os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	content := string(requestBody)
	headerName := "ej_" + strings.ReplaceAll(ip_address, ".", "_")
	filename := storePath + "/" + headerName

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(content + "\n"); err != nil {
		return err
	}

	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	return nil
}

func ConsumeFileEjol() error {
	date := time.Now().Format("20060102")
	dirPath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + date + "/"
	sectionDirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Printf("Error Get Dir : %s", err)
		return err
	}

	for _, files := range sectionDirs {
		fmt.Println(files.Name())

		ej_content, err := os.ReadFile(dirPath + files.Name())
		if err != nil {
			models.ErrorLogger.Printf("Error read file : %s", err)
			return err
		}

		ip_address := strings.ReplaceAll(files.Name(), "_", ".")
		ip_address = ip_address[3:]
		// fmt.Println(ip_address)
		getKanwil, found := Cac.Get(ip_address)
		if !found {
			log.Printf("RC : %d - Cache Ip : %s Not Found \n", http.StatusNotFound, ip_address)
			continue
		}
		kanwil := getKanwil.(string)

		err = processRequestEjlog(string(ej_content), ip_address, kanwil)
		if err != nil {
			models.ErrorLogger.Printf("Error processRequestEjlog : %s", err)
			return err
		}
		models.InfoLogger.Printf("Sukses menyimpan RequestEjlog dalam file %s", dirPath+files.Name())

	}

	return nil
}

func ConsumeFileEjolSchedule() error {
	date := time.Now().Format("20060102")
	dirPath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + date + "/"
	sectionDirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Printf("Error Get Dir : %s", err)
		return err
	}
	cache_directory := time.Now().Format("20060102") + "/"
	var dst []byte
	expire := 25 * time.Hour

	for _, files := range sectionDirs {
		ip_address := strings.ReplaceAll(files.Name(), "_", ".")
		ip_address = ip_address[3:]
		// fmt.Println(ip_address)
		getKanwil, found := Cac.Get(ip_address)
		if !found {
			models.ErrorLogger.Printf("RC : %d - Cache Ip : %s Not Found \n", http.StatusNotFound, ip_address)
			continue
		}
		kanwil := getKanwil.(string)

		_, err = filecache.Get(cache_directory, strings.ReplaceAll(ip_address, ".", "_"), &dst)
		if err != nil {
			models.ErrorLogger.Printf("IP Not Found in cache %s - Error : %s \n", ip_address, err)
			return err
		}
		cachevalue, _ := strconv.Atoi(string(dst))

		tblname := strings.ReplaceAll(ip_address, ".", "_")
		fn := dirPath + "ej_" + tblname

		line, lastl, err := ReadLineStreamUntilFinish(fn, cachevalue, 0)
		if err != nil {
			models.ErrorLogger.Printf("Error %s \n", err)
			return err
		}

		err = processRequestEjlogSchedule(line, ip_address, kanwil)
		if err != nil {
			models.ErrorLogger.Printf("Error processRequestEjlog : %s \n", err)
			return err
		}

		lastlinestring := strconv.Itoa(lastl + 1)
		data_lastline := []byte(lastlinestring)
		err_cache := filecache.Set(cache_directory, strings.ReplaceAll(ip_address, ".", "_"), data_lastline, expire)
		if err_cache != nil {
			models.InfoLogger.Printf("Failed to set new value cache for IP %s - Error : %s \n", ip_address, err)
			return err
		}
		models.InfoLogger.Printf("Data Saved Successfully %s \n", dirPath+files.Name())

	}

	return nil
}

func processRequestEjlog(ejcontent, ip_address, kanwil string) error {
	tbl_name := "ej_" + strings.ReplaceAll(ip_address, ".", "_")
	tbl_withdraw := "ejlog_withdraw"
	tbl_print_cash := "ejlog_print_cash"
	tbl_emergency_receipt := "ejlog_emergency_receipt"
	tbl_addcash := "ejlog_add_cash"

	dbname := "ejlog_" + kanwil + "_" + time.Now().Format("20060102") //make format dbname : ejlog_kanwil_yyyymmdd
	requestBody := ejcontent
	ejol_map := strings.Split(string(requestBody), "\n")
	db, err := config.DbConn(dbname)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	//insert data
	for _, ejlog_line := range ejol_map {
		ejlog_line_tmp := strings.ReplaceAll(strings.ToUpper(ejlog_line), " ", "")
		if ejlog_line_tmp == "" {
			continue
		}
		result, err := tx.Exec("INSERT INTO "+tbl_name+" (ip_address, ejlog) VALUES(?, ?)", ip_address, ejlog_line)
		if err != nil {
			// tx.Rollback()
			return err
		}
		tblej_id, _ := result.LastInsertId() //getLastIndex Executed

		//WITHDRAW
		if strings.Contains(ejlog_line_tmp, "SWITCHING") {
			_, err := tx.Exec("INSERT INTO `"+tbl_withdraw+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//PRINT CASH
		if strings.Contains(ejlog_line_tmp, "PRINTCASH") ||
			strings.Contains(ejlog_line_tmp, "IDR,IDR,IDR,IDR") ||
			strings.Contains(ejlog_line_tmp, "TYPE1TYPE2") ||
			strings.Contains(ejlog_line_tmp, "CASHTOTALTYPE1TYPE2") ||
			// strings.Contains(ejlog_line_tmp, "[020t[05pTYPE1TYPE2") ||
			strings.Contains(ejlog_line_tmp, "[020t\u001B[05pTYPE1TYPE2") ||
			strings.Contains(ejlog_line_tmp, "#CURDENOCST+REJ=REM+DISP=TOTAL") ||
			strings.Contains(ejlog_line_tmp, "CURDENOINITDISPDEPCSTRJ") ||
			strings.Contains(ejlog_line_tmp, "CASHCOUNTINFO") {
			_, err := tx.Exec("INSERT INTO `"+tbl_print_cash+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//EMERGENCY RECEIPT
		if strings.Contains(ejlog_line_tmp, "EMERGENCYRECEIPT") ||
			strings.Contains(ejlog_line_tmp, "TRANSACTIONJRNLTRANSACTION") ||
			strings.Contains(ejlog_line_tmp, "PLEASECONTACTBRANCH") {
			_, err := tx.Exec("INSERT INTO `"+tbl_emergency_receipt+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//ADD CASH
		if strings.Contains(ejlog_line_tmp, "ADDCASH") ||
			strings.Contains(ejlog_line_tmp, "CASHCOUNTERS") ||
			strings.Contains(ejlog_line_tmp, "[05pCASHADDED") ||
			strings.Contains(ejlog_line_tmp, "CLEARCASH") ||
			strings.Contains(ejlog_line_tmp, "REPLENISHMENT") {
			_, err := tx.Exec("INSERT INTO `"+tbl_addcash+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	defer db.Close()

	return nil
}

func processRequestEjlog2(ejcontent, ip_address, kanwil string) (string, error) {
	tbl_name := "ej_" + strings.ReplaceAll(ip_address, ".", "_")
	tbl_withdraw := "ejlog_withdraw"
	tbl_print_cash := "ejlog_print_cash"
	tbl_emergency_receipt := "ejlog_emergency_receipt"
	tbl_addcash := "ejlog_add_cash"

	dbname := "ejlog_" + kanwil + "_" + time.Now().Format("20060102") //make format dbname : ejlog_kanwil_yyyymmdd
	requestBody := ejcontent
	ejol_map := strings.Split(string(requestBody), "\n")
	db, err := config.DbConn(dbname)
	if err != nil {
		return ejol_map[0], err
	}

	tx, err := db.Begin()
	if err != nil {
		return ejol_map[0], err
	}

	//insert data
	for _, ejlog_line := range ejol_map {
		ejlog_line_tmp := strings.ReplaceAll(strings.ToUpper(ejlog_line), " ", "")
		if ejlog_line_tmp == "" {
			continue
		}
		result, err := tx.Exec("INSERT INTO "+tbl_name+" (ip_address, ejlog) VALUES(?, ?)", ip_address, ejlog_line)
		if err != nil {
			// tx.Rollback()
			return ejlog_line, err
		}
		tblej_id, _ := result.LastInsertId() //getLastIndex Executed

		//WITHDRAW
		if strings.Contains(ejlog_line_tmp, "SWITCHING") {
			_, err := tx.Exec("INSERT INTO `"+tbl_withdraw+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return ejlog_line, err
			}
		}

		//PRINT CASH
		if strings.Contains(ejlog_line_tmp, "PRINTCASH") ||
			strings.Contains(ejlog_line_tmp, "IDR,IDR,IDR,IDR") ||
			strings.Contains(ejlog_line_tmp, "TYPE1TYPE2") ||
			strings.Contains(ejlog_line_tmp, "CASHTOTALTYPE1TYPE2") ||
			// strings.Contains(ejlog_line_tmp, "[020t[05pTYPE1TYPE2") ||
			strings.Contains(ejlog_line_tmp, "[020t\u001B[05pTYPE1TYPE2") ||
			strings.Contains(ejlog_line_tmp, "#CURDENOCST+REJ=REM+DISP=TOTAL") ||
			strings.Contains(ejlog_line_tmp, "CURDENOINITDISPDEPCSTRJ") ||
			strings.Contains(ejlog_line_tmp, "CASHCOUNTINFO") {
			_, err := tx.Exec("INSERT INTO `"+tbl_print_cash+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return ejlog_line, err
			}
		}

		//EMERGENCY RECEIPT
		if strings.Contains(ejlog_line_tmp, "EMERGENCYRECEIPT") ||
			strings.Contains(ejlog_line_tmp, "TRANSACTIONJRNLTRANSACTION") ||
			strings.Contains(ejlog_line_tmp, "PLEASECONTACTBRANCH") {
			_, err := tx.Exec("INSERT INTO `"+tbl_emergency_receipt+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return ejlog_line, err
			}
		}

		//ADD CASH
		if strings.Contains(ejlog_line_tmp, "ADDCASH") ||
			strings.Contains(ejlog_line_tmp, "CASHCOUNTERS") ||
			strings.Contains(ejlog_line_tmp, "[05pCASHADDED") ||
			strings.Contains(ejlog_line_tmp, "CLEARCASH") ||
			strings.Contains(ejlog_line_tmp, "REPLENISHMENT") {
			_, err := tx.Exec("INSERT INTO `"+tbl_addcash+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return ejlog_line, err
			}
		}
	}
	tx.Commit()
	defer db.Close()
	return ejol_map[0], nil
}

func processRequestEjlogSchedule(ejcontent []string, ip_address string, kanwil string) error {
	tbl_name := "ej_" + strings.ReplaceAll(ip_address, ".", "_")
	tbl_withdraw := "ejlog_withdraw"
	tbl_print_cash := "ejlog_print_cash"
	tbl_emergency_receipt := "ejlog_emergency_receipt"
	tbl_addcash := "ejlog_add_cash"

	dbname := "ejlog_" + kanwil + "_" + time.Now().Format("20060102") //make format dbname : ejlog_kanwil_yyyymmdd
	// requestBody := ejcontent
	// ejol_map := strings.Split(requestBody, "\n")
	ejol_map := ejcontent
	db, err := config.DbConn(dbname)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	fmt.Println(len(ejol_map))
	//insert data
	for i := 0; i < len(ejol_map); i++ {
		ejlog_line := ejol_map[i]
		ejlog_line_tmp := strings.ReplaceAll(strings.ToUpper(ejlog_line), " ", "")
		if ejlog_line_tmp == "" {
			continue
		}
		result, err := tx.Exec("INSERT INTO "+tbl_name+" (ip_address, ejlog) VALUES(?, ?)", ip_address, ejlog_line)
		if err != nil {
			// tx.Rollback()
			return err
		}
		tblej_id, _ := result.LastInsertId() //getLastIndex Executed

		//WITHDRAW
		if strings.Contains(ejlog_line_tmp, "SWITCHING") {
			_, err := tx.Exec("INSERT INTO `"+tbl_withdraw+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//PRINT CASH
		if strings.Contains(ejlog_line_tmp, "PRINTCASH") ||
			strings.Contains(ejlog_line_tmp, "IDR,IDR,IDR,IDR") ||
			strings.Contains(ejlog_line_tmp, "TYPE1TYPE2") ||
			strings.Contains(ejlog_line_tmp, "CASHTOTALTYPE1TYPE2") ||
			// strings.Contains(ejlog_line_tmp, "[020t[05pTYPE1TYPE2") ||
			strings.Contains(ejlog_line_tmp, "[020t\u001B[05pTYPE1TYPE2") ||
			strings.Contains(ejlog_line_tmp, "#CURDENOCST+REJ=REM+DISP=TOTAL") ||
			strings.Contains(ejlog_line_tmp, "CURDENOINITDISPDEPCSTRJ") ||
			strings.Contains(ejlog_line_tmp, "CASHCOUNTINFO") {
			_, err := tx.Exec("INSERT INTO `"+tbl_print_cash+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//EMERGENCY RECEIPT
		if strings.Contains(ejlog_line_tmp, "EMERGENCYRECEIPT") ||
			strings.Contains(ejlog_line_tmp, "TRANSACTIONJRNLTRANSACTION") ||
			strings.Contains(ejlog_line_tmp, "PLEASECONTACTBRANCH") {
			_, err := tx.Exec("INSERT INTO `"+tbl_emergency_receipt+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//ADD CASH
		if strings.Contains(ejlog_line_tmp, "ADDCASH") ||
			strings.Contains(ejlog_line_tmp, "CASHCOUNTERS") ||
			strings.Contains(ejlog_line_tmp, "[05pCASHADDED") ||
			strings.Contains(ejlog_line_tmp, "CLEARCASH") ||
			strings.Contains(ejlog_line_tmp, "REPLENISHMENT") {
			_, err := tx.Exec("INSERT INTO `"+tbl_addcash+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, ejlog_line, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		// break
	}
	tx.Commit()
	defer db.Close()

	return nil
}

func ReadLineStreamUntilFinish(fn string, lineNum int, lineLength int) (line []string, lastLine int, err error) {
	// func ReadLineStreamUntilFinish(fn string, lineNum int, lineLength int) (line string, lastLine int, err error) {
	// line_length := 30
	// line_length := 0
	// if lineLength > 0 {
	// 	line_length = lineLength
	// }
	file, err := os.Open(fn)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var txtlines []string
	var line_found bool
	line_found = false
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lastLine++ //start dari 0

		if lastLine == lineNum {
			line_found = true
		}

		// if line_length == 0 {
		// 	txtlines = append(txtlines, sc.Text())
		// 	return txtlines, lastLine, sc.Err()
		// }

		// txtlines = append(txtlines, sc.Text())
		if line_found {
			txtlines = append(txtlines, sc.Text())
			// line_length--
		}

	}

	if lastLine < lineNum {
		fmt.Printf("lastLine %d < lineNum %d return EOF \n", lastLine, lineNum)

		return nil, lastLine, io.EOF
	}
	// fmt.Printf("Last line %d \n", lastLine)
	// fmt.Printf("Content %s \n", strings.Join(txtlines, "\n"))

	// return strings.Join(txtlines, "\n"), lastLine, nil
	return txtlines, lastLine, nil
}

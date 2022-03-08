//NXLOG VERSI 3 MENYIMPAN DALAM BENTUK FILE
package controller

import (
	"ejol/ejlog-server/config"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var Cac *cache.Cache

func V3MultilineWincor_1(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	start := time.Now()

	todayDate := start.Format("20060102")
	ip_address, _, error := net.SplitHostPort(r.RemoteAddr)

	if error != nil {
		ErrorLogger.Printf("RC : %d - Error : %s", http.StatusNotAcceptable, error)
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
	tblname := strings.ReplaceAll(ip_address, ".", "_")

	// storePath := "appendrow/" + todayDate + "/" + tblname
	storePath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate + "/" + tblname

	if _, errPath := os.Stat(storePath); os.IsNotExist(errPath) {
		err := os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			ErrorLogger.Printf("Error : %s", err)
			return
		}
	}

	rb, _ := ioutil.ReadAll(r.Body)
	content := string(rb)
	headerName := ip_address + "_" + strings.ReplaceAll(start.Format("150405.0000000"), ".", "")
	filename := storePath + "/" + headerName
	err := ioutil.WriteFile(filename, []byte(content), 0666)
	if err != nil {
		ErrorLogger.Printf("RC : %d - Error WriteFile : %s", http.StatusNotAcceptable, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = processRequestEjlog2(content, ip_address, kanwil)
	if err != nil {
		ErrorLogger.Printf("RC : %d - Error processRequestEjlog2 : %s", http.StatusNotAcceptable, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	InfoLogger.Printf("RC : %d - Save successfully.", http.StatusOK)
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

	fmt.Printf("Ip address %s \n", ip_address)

	getKanwil, found := Cac.Get(ip_address)
	if !found {
		ErrorLogger.Printf("RC : %d - Ip : %s Not Found ", http.StatusNotFound, ip_address)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	kanwil := getKanwil.(string)
	tblname := strings.ReplaceAll(ip_address, ".", "_")

	// storePath := "appendrow/" + todayDate + "/" + tblname
	storePath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate + "/" + tblname
	rb, _ := ioutil.ReadAll(r.Body)
	content := string(rb)
	// headerName := ip_address + "_" + strings.ReplaceAll(start.Format("150405.0000000"), ".", "")
	headerName := ip_address
	filename := storePath + "/" + headerName

	if _, errPath := os.Stat(storePath); os.IsNotExist(errPath) {
		err := os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			ErrorLogger.Printf("Ip : %s Error : %s", ip_address, err)
			return
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ErrorLogger.Printf("Ip : %s Error : %s", ip_address, err)
		return
	}
	defer file.Close()
	if _, err := file.WriteString(content + "\n"); err != nil {
		ErrorLogger.Fatal(err)
		return
	}

	err = processRequestEjlog2(content, ip_address, kanwil)
	if err != nil {
		ErrorLogger.Printf("RC : %d - %s Error processRequestEjlog2 : %s", http.StatusNotAcceptable, ip_address, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	InfoLogger.Printf("RC : %d - %s", http.StatusOK, ip_address)
	w.WriteHeader(http.StatusOK)
}

func V3MultilineWincorAppendFile(w http.ResponseWriter, r *http.Request) {
	// ip_address, _, error := net.SplitHostPort(r.RemoteAddr)
	// if error != nil {
	// 	ErrorLogger.Printf("RC : %d - %s Error : %s", http.StatusNotAcceptable, ip_address, error)
	// 	w.WriteHeader(http.StatusNotAcceptable)
	// 	return
	// }
	ip_address := r.Header.Get("IP-ADDRESS")

	if ip_address == "::1" {
		ip_address = "127.0.0.1"
	}

	kanwil, found := Cac.Get(ip_address)
	if !found {
		ErrorLogger.Printf("RC : %d - Ip %s Not Found ", http.StatusNotFound, ip_address)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	requestBody, _ := ioutil.ReadAll(r.Body)
	ejol_map := strings.Split(string(requestBody), "\n")
	namefile := ejol_map[0]

	if strings.Contains(namefile, strings.ToUpper("TRANSFER_OUTPUT")) {
		fmt.Println("inCURl")
		err := AgentByCURL(ip_address, string(requestBody), ejol_map[0], kanwil.(string))
		if err != nil {
			ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error %s", err)
			return
		}
	} else {
		fmt.Println("inAgent")
		err := AgentByNxLogAppendFile(ip_address, string(requestBody), kanwil.(string))
		if err != nil {
			ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error %s", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	InfoLogger.Printf("RC : %d - %s", http.StatusOK, ip_address)
	fmt.Fprintf(w, "200")
}

func V3MultilineWincorSplitFile(w http.ResponseWriter, r *http.Request) {
	ip_address, _, error := net.SplitHostPort(r.RemoteAddr)
	if error != nil {
		ErrorLogger.Printf("RC : %d - %s Error : %s", http.StatusNotAcceptable, ip_address, error)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if ip_address == "::1" {
		ip_address = "127.0.0.1"
	}

	kanwil, found := Cac.Get(ip_address)
	if !found {
		ErrorLogger.Printf("RC : %d - Ip %s Not Found ", http.StatusNotFound, ip_address)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	requestBody, _ := ioutil.ReadAll(r.Body)
	ejol_map := strings.Split(string(requestBody), "\n")
	namefile := ejol_map[0]

	if strings.EqualFold(namefile[0:15], strings.ToUpper("TRANSFER_OUTPUT")) {
		err := AgentByCURL(ip_address, string(requestBody), ejol_map[0], kanwil.(string))
		if err != nil {
			ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error %s", err)
			return
		}
	} else {
		err := AgentByNxLog(ip_address, string(requestBody), kanwil.(string))
		if err != nil {
			ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error %s", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	InfoLogger.Printf("RC : %d - %s", http.StatusOK, ip_address)
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
			ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
			return err
		}
	}

	content := string(requestBody)
	filesave := storePath + "/" + filename
	err := ioutil.WriteFile(filesave, []byte(content), 0666)
	if err != nil {
		ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
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
			ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
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
		ErrorLogger.Printf("[RC : %d - %s] Error : %s", http.StatusOK, ip_address, err)
		return err
	}

	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	return nil
}

func AgentByNxLogAppendFile(ip_address, requestBody, kanwil string) error {
	start := time.Now()
	todayDate := start.Format("20060102")

	tblname := strings.ReplaceAll(ip_address, ".", "_")
	storePath := os.Getenv("EJOL_DIRECTORY_FILE") + "appendrow/" + todayDate + "/" + tblname
	// storePath := "appendrow/" + todayDate + "/" + tblname

	if _, errPath := os.Stat(storePath); os.IsNotExist(errPath) {
		err := os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			ErrorLogger.Printf("Ip : %s Error : %s", ip_address, err)
			return err
		}
	}

	content := string(requestBody)
	headerName := ip_address
	filename := storePath + "/" + headerName

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ErrorLogger.Printf("Ip : %s Error : %s", ip_address, err)
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(content + "\n"); err != nil {
		ErrorLogger.Fatal(err)
		return err
	}

	end := <-time.After(10 * time.Millisecond)
	end.Sub(start)
	return nil
}

func ConsumeFileEjol() error {
	date := time.Now().Format("20060102")
	// dirPath := os.Getenv("EJOL_DIRECTORY") + "project/golang/go-ejlog-server/appendrow/" + date + "/127_0_0_1/"
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
			err = processRequestEjlog(string(ej_content), ip_address, kanwil)
			if err != nil {
				ErrorLogger.Printf("Error processRequestEjlog : %s", err)
				return err
			}
			InfoLogger.Printf("Sukses menyimpan RequestEjlog dalam file %s", filePath+f.Name())

			err = os.Rename(filePath+f.Name(), filePath+"D"+f.Name())
			if err != nil {
				ErrorLogger.Printf("Error pada merubah file : %s", err)
				return err
			}

			InfoLogger.Printf("Sukses rename File %s", filePath+f.Name())
		}

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
	for _, element := range ejol_map {
		result, err := tx.Exec("INSERT INTO "+tbl_name+" (ip_address, ejlog) VALUES(?, ?)", ip_address, element)
		if err != nil {
			// tx.Rollback()
			return err
		}
		tblej_id, _ := result.LastInsertId() //getLastIndex Executed

		//WITHDRAW
		if strings.Contains(element, "SWITCHING") {
			_, err := tx.Exec("INSERT INTO `"+tbl_withdraw+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, element, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//PRINT CASH
		if strings.Contains(element, "PRINTCASH") ||
			strings.Contains(element, "IDR,IDR,IDR,IDR") ||
			strings.Contains(element, "TYPE1TYPE2") ||
			strings.Contains(element, "CASHTOTALTYPE1TYPE2") ||
			strings.Contains(element, "[020t[05pTYPE1TYPE2") ||
			strings.Contains(element, "#CURDENOCST+REJ=REM+DISP=TOTAL") ||
			strings.Contains(element, "CURDENOINITDISPDEPCSTRJ") ||
			strings.Contains(element, "CASHCOUNTINFO") {
			_, err := tx.Exec("INSERT INTO `"+tbl_print_cash+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, element, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//EMERGENCY RECEIPT
		if strings.Contains(element, "EMERGENCYRECEIPT") ||
			strings.Contains(element, "TRANSACTIONJRNLTRANSACTION") ||
			strings.Contains(element, "PLEASECONTACTBRANCH") {
			_, err := tx.Exec("INSERT INTO `"+tbl_emergency_receipt+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, element, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//ADD CASH
		if strings.Contains(element, "ADDCASH") ||
			strings.Contains(element, "CASHCOUNTERS") ||
			strings.Contains(element, "[05pCASHADDED") ||
			strings.Contains(element, "CLEARCASH") ||
			strings.Contains(element, "REPLENISHMENT") {
			_, err := tx.Exec("INSERT INTO `"+tbl_addcash+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, element, 0)
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

func processRequestEjlog2(ejcontent, ip_address, kanwil string) error {
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
	for _, element := range ejol_map {
		result, err := tx.Exec("INSERT INTO "+tbl_name+" (ip_address, ejlog) VALUES(?, ?)", ip_address, element)
		if err != nil {
			// tx.Rollback()
			return err
		}
		tblej_id, _ := result.LastInsertId() //getLastIndex Executed

		//WITHDRAW
		if strings.Contains(element, "SWITCHING") {
			_, err := tx.Exec("INSERT INTO `"+tbl_withdraw+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, element, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//PRINT CASH
		if strings.Contains(element, "PRINTCASH") ||
			strings.Contains(element, "IDR,IDR,IDR,IDR") ||
			strings.Contains(element, "TYPE1TYPE2") ||
			strings.Contains(element, "CASHTOTALTYPE1TYPE2") ||
			strings.Contains(element, "[020t[05pTYPE1TYPE2") ||
			strings.Contains(element, "#CURDENOCST+REJ=REM+DISP=TOTAL") ||
			strings.Contains(element, "CURDENOINITDISPDEPCSTRJ") ||
			strings.Contains(element, "CASHCOUNTINFO") {
			_, err := tx.Exec("INSERT INTO `"+tbl_print_cash+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, element, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//EMERGENCY RECEIPT
		if strings.Contains(element, "EMERGENCYRECEIPT") ||
			strings.Contains(element, "TRANSACTIONJRNLTRANSACTION") ||
			strings.Contains(element, "PLEASECONTACTBRANCH") {
			_, err := tx.Exec("INSERT INTO `"+tbl_emergency_receipt+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, element, 0)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//ADD CASH
		if strings.Contains(element, "ADDCASH") ||
			strings.Contains(element, "CASHCOUNTERS") ||
			strings.Contains(element, "[05pCASHADDED") ||
			strings.Contains(element, "CLEARCASH") ||
			strings.Contains(element, "REPLENISHMENT") {
			_, err := tx.Exec("INSERT INTO `"+tbl_addcash+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, element, 0)
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

package controller

import (
	"context"
	"ejol/ejlog-server/config"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	fmt.Println("access here first")
	err := godotenv.Load(".env")
	if err != nil {
		ErrorLogger.Fatal("Error load file env : ", err)
	}
	filename := os.Getenv("EJOL_DIRECTORY_LOG") + "ejlog-server-" + time.Now().Format("20060102") + ".log"
	// filename := "ejlog-server-" + time.Now().Format("20060102") + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)
}

func MultilineWincor(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
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
	tbl_name := "ej_" + strings.ReplaceAll(ip_address, ".", "_")
	tbl_withdraw := "ejlog_withdraw"
	tbl_print_cash := "ejlog_print_cash"
	tbl_emergency_receipt := "ejlog_emergency_receipt"
	tbl_addcash := "ejlog_add_cash"
	dbname := "ejlog_" + kanwil + "_" + time.Now().Format("20060102") //make format dbname : ejlog_kanwil_yyyymmdd

	requestBody, _ := ioutil.ReadAll(r.Body)
	ejol_map := strings.Split(string(requestBody), "\n")
	db, err := config.DbConn(dbname)
	if err != nil {
		ErrorLogger.Printf("RC : %d - Error connect to DB : %s", http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ErrorLogger.Printf("RC : %d - Error begin transaction %s", http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//insert data
	for _, element := range ejol_map {
		result, err := tx.Exec("INSERT INTO "+tbl_name+" (ip_address, ejlog) VALUES(?, ?)", ip_address, element)
		if err != nil {
			ErrorLogger.Printf("RC : %d - Error : %s", http.StatusInternalServerError, err)
			w.WriteHeader(http.StatusInternalServerError)
			db.Close()
			return
		}
		tblej_id, _ := result.LastInsertId() //getLastIndex Executed

		//WITHDRAW
		if strings.Contains(element, "SWITCHING") {
			_, err := tx.Exec("INSERT INTO `"+tbl_withdraw+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, element, 0)
			if err != nil {
				ErrorLogger.Printf("RC : %d - Error WD: %s", http.StatusInternalServerError, err)
				w.WriteHeader(http.StatusInternalServerError)
				tx.Rollback()
				db.Close()
				return
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
				ErrorLogger.Printf("RC : %d - Error PrintCash: %s", http.StatusInternalServerError, err)
				w.WriteHeader(http.StatusInternalServerError)
				tx.Rollback()
				db.Close()
				return
			}
		}

		//EMERGENCY RECEIPT
		if strings.Contains(element, "EMERGENCYRECEIPT") ||
			strings.Contains(element, "TRANSACTIONJRNLTRANSACTION") ||
			strings.Contains(element, "PLEASECONTACTBRANCH") {
			_, err := tx.Exec("INSERT INTO `"+tbl_emergency_receipt+"`(`index`, `dbname`, `ip_address`, `ejlog`, `is_read`) VALUES (?, ?, ?, ?, ?)", tblej_id, tbl_name, ip_address, element, 0)
			if err != nil {
				ErrorLogger.Printf("RC : %d - Error Emergency: %s", http.StatusInternalServerError, err)
				w.WriteHeader(http.StatusInternalServerError)
				tx.Rollback()
				db.Close()
				return
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
				ErrorLogger.Printf("RC : %d - Error Addcash: %s", http.StatusInternalServerError, err)
				w.WriteHeader(http.StatusInternalServerError)
				tx.Rollback()
				db.Close()
				return
			}
		}
	}
	tx.Commit()
	db.Close()
	// defer db.Close()

	elapsed := time.Since(start)
	InfoLogger.Printf("RC : %d - This request took %s ", http.StatusOK, elapsed)
	w.WriteHeader(http.StatusOK)
}

func checkTable(tbl_name, kanwil string) error {
	dbname := "ejlog_" + kanwil + "_" + time.Now().Format("20060102")

	db, err := config.DbConn(dbname)
	if err != nil {
		ErrorLogger.Printf("Error connection to db : ", err)
		return err
	}

	query := `CREATE TABLE IF NOT EXISTS ` + tbl_name + `(
		id INT(10) NOT NULL AUTO_INCREMENT,
		ip_address VARCHAR(50) NOT NULL,
		row_ej VARCHAR(10) NULL DEFAULT NULL,
		ejlog LONGTEXT NOT NULL,
		header_ej VARCHAR(100) NULL DEFAULT NULL,
		is_curl VARCHAR(2) NOT NULL DEFAULT '1',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL DEFAULT NULL,
		PRIMARY KEY (id) USING BTREE,
		INDEX created_at (created_at) USING BTREE
		)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err2 := db.ExecContext(ctx, query)
	if err2 != nil {
		ErrorLogger.Printf("Error %s when creating ejlog table - ", err2)
		return err2
	}
	defer db.Close()

	return nil
}

func checkIpValid(ip_address string) (string, error) {
	var tid, ip_atm, _type, kanwil2 string
	db, err := config.DbConn("ejlog3")
	if err != nil {
		ErrorLogger.Printf("Error connection to db : ", err)
		return "", err
	}
	err2 := db.QueryRow("select tid, ip_address, type, kanwil2 from atm_mappings where ip_address = ? limit 1", ip_address).Scan(&tid, &ip_atm, &_type, &kanwil2)
	if err2 != nil {
		return "", err2
	}
	defer db.Close()

	return kanwil2, nil
}

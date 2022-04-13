package job

import (
	"ejol/ejlog-server/config"
	"ejol/ejlog-server/controller"
	"ejol/ejlog-server/models"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/miguelmota/go-filecache"
	"github.com/patrickmn/go-cache"
)

func JobExportCountAtm() {
	for {
		fmt.Println("Call Scheduler JobExportCountAtm")
		db, err := config.DbConn("ejlog3")
		if err != nil {
			controller.ErrorLogger.Println("Error connect to DB : ", err)
			return
		}
		defer db.Close()

		atmMappingModel := models.AtmMappingCacheModel{
			DB: db,
		}

		total_atm, err := atmMappingModel.CountAtmMapping()
		if err != nil {
			controller.ErrorLogger.Printf("Error model : %s", err)
			return
		}

		os.Setenv("TOTAL_MESIN", strconv.Itoa(total_atm))
		time.Sleep(24 * time.Hour)
		controller.InfoLogger.Println("After 24 Hour")
	}
}

func JobCacheAtmMappings() {
	start := time.Now()
	for {
		fmt.Println("Call Scheduler")

		// c := cache.New(25*time.Hour, 26*time.Hour)
		controller.Cac = cache.New(25*time.Hour, 26*time.Hour)
		controller.Cac.Flush()

		db, err := config.DbConn("ejlog3")
		if err != nil {
			controller.ErrorLogger.Println("Error connect to DB : ", err)
			return
		}
		// defer db.Close()

		atmMappingModel := models.AtmMappingCacheModel{
			DB: db,
		}

		atms, err2 := atmMappingModel.GetData()
		if err2 != nil {
			controller.InfoLogger.Printf("Error get atms : %s", err2)
		}
		for _, atm := range atms {
			controller.Cac.Set(atm.IpAddress, atm.Kanwil2, cache.NoExpiration)
		}

		elapsed := time.Since(start)
		fmt.Printf("%s tooks %s\n", "Method", elapsed)
		time.Sleep(24 * time.Hour)
		controller.InfoLogger.Println("After 10 Minutes")
	}

}

func JobCacheAtmMappings2() {
	for {
		fmt.Println("Call Scheduler")

		c := cache.New(25*time.Hour, 26*time.Hour)
		c.Flush()

		db, err := config.DbConn("ejlog3")
		if err != nil {
			controller.ErrorLogger.Println("Error connect to DB : ", err)
			return
		}
		defer db.Close()

		rows, err := db.Query("select ip_address, kanwil2 from atm_mappings limit ?,?", 0, 5000)
		var i = 0
		if err != nil {
			controller.ErrorLogger.Println("Error query : ", err)
		} else {
			for rows.Next() {
				var ipaddr, kanwil2 string
				i++
				err2 := rows.Scan(&ipaddr, &kanwil2)
				if err2 != nil {
					controller.ErrorLogger.Println("Error looping data : ", err2)
				}

				c.Set(ipaddr, kanwil2, cache.NoExpiration)
				controller.InfoLogger.Printf("Data %d - Ip Addr %s :  - Kanwil : %s", i, ipaddr, kanwil2)

			}

			if err := rows.Err(); err != nil {
				panic(err)
			}
		}

		time.Sleep(10 * time.Minute)
		controller.InfoLogger.Println("After 10 Minutes")
	}
}

func TestingCache3() {
	start := time.Now()
	expire := 25 * time.Hour

	db, err := config.DbConn("ejlog3")
	if err != nil {
		fmt.Printf("Error : %s \n", err)
	}

	atmMappingModel := models.AtmMappingModel{
		DB: db,
	}

	cache_directory := time.Now().Format("20060102") + "/"
	machines, err := atmMappingModel.GetDataKanwilCache()
	if err != nil {
		fmt.Printf("Error : %s \n", err)
	}

	for _, atm := range machines {
		err := filecache.Set(cache_directory, strings.ReplaceAll(atm.IpAddress, ".", "_"), []byte("1"), expire)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("set cache for %s \n", atm.IpAddress)
	}

	fmt.Println("Job Finished.")
	elapsed := time.Since(start)
	fmt.Printf("This process took %s \n", elapsed)
}

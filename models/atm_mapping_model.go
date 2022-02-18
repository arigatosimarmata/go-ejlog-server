package models

import (
	"database/sql"
	"ejol/ejlog-server/controller"
)

type AtmMappingCache struct {
	IpAddress string
	Kanwil2   string
}

type AtmMappingCacheModel struct {
	DB *sql.DB
}

func (atmModel AtmMappingCacheModel) GetData() ([]AtmMappingCache, error) {
	rows, err := atmModel.DB.Query("select ip_address, kanwil2 from atm_mappings")
	var i = 0
	if err != nil {
		controller.ErrorLogger.Println("Error query : ", err)
		return nil, err
	} else {
		atms := []AtmMappingCache{}
		for rows.Next() {
			var ipaddr, kanwil2 string
			i++
			err2 := rows.Scan(&ipaddr, &kanwil2)
			if err2 != nil {
				controller.ErrorLogger.Println("Error looping data : ", err2)
				return nil, err2
			} else {
				controller.InfoLogger.Printf("Data %d - Ip Addr %s :  - Kanwil : %s", i, ipaddr, kanwil2)
				atm := AtmMappingCache{ipaddr, kanwil2}
				atms = append(atms, atm)
			}
		}
		atmModel.DB.Close()
		return atms, nil
	}
}

func (atmModel AtmMappingCacheModel) Limit(offset, count int) ([]AtmMappingCache, error) {
	rows, err := atmModel.DB.Query("select ip_address, kanwil2 from atm_mappings limit ?,?", offset, count)
	var i = 0
	if err != nil {
		controller.ErrorLogger.Println("Error query : ", err)
		return nil, err
	} else {
		atms := []AtmMappingCache{}
		for rows.Next() {
			var ipaddr, kanwil2 string
			i++
			err2 := rows.Scan(&ipaddr, &kanwil2)
			if err2 != nil {
				controller.ErrorLogger.Println("Error looping data : ", err2)
				return nil, err2
			} else {
				controller.InfoLogger.Printf("Data %d - Ip Addr %s :  - Kanwil : %s", i, ipaddr, kanwil2)
				atm := AtmMappingCache{ipaddr, kanwil2}
				atms = append(atms, atm)
			}
		}
		defer atmModel.DB.Close()
		return atms, nil
	}
}

func (atmModel AtmMappingCacheModel) GetDataV2() ([]AtmMappingCache, error) {
	records := make([]AtmMappingCache, 0)
	count := 24000
	bucketSize := 5000
	resultCount := 0
	resultChannel := make(chan []AtmMappingCache, 0)
	for beginID := 1; beginID <= count; beginID += bucketSize {
		endId := beginID + bucketSize
		go func(beginID int, endId int) {
			var ipaddr, kanwil2 string
			currentRecords := make([]AtmMappingCache, 0)
			err := atmModel.DB.QueryRow("select ip_address, kanwil2 from atm_mappings").Scan(&ipaddr, &kanwil2)
			// err := atmModel.DB.QueryRow("select ip_address, kanwil2 from atm_mappings").Scan(&AtmMappingCache)
			if err != nil {
				controller.ErrorLogger.Printf("Error query %s : ", err)
			}
			controller.InfoLogger.Printf("Ip Address : %s - Kanwil : %s", ipaddr, kanwil2)
			resultChannel <- currentRecords
		}(beginID, endId)
		resultCount += 1
	}

	for i := 0; i < resultCount; i++ {
		currentRecords := <-resultChannel
		records = append(records, currentRecords...)
	}

	return records, nil
}

func (atmModel AtmMappingCacheModel) CountAtmMapping() (int, error) {
	var total_atmmapping int
	err := atmModel.DB.QueryRow("select count(1) as total_row from atm_mappings").Scan(&total_atmmapping)
	if err != nil {
		return 0, err
	}

	return total_atmmapping, nil
}

package controller

import (
	"ejol/ejlog-server/models"
	"encoding/json"
	"os"
)

func LoadConfigKeyword(filename string) (map[string]map[string]string, error) {
	var v map[string]map[string]string
	configFile, err := os.Open(filename)
	if err != nil {
		models.ErrorLogger.Println(err)
		return nil, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&v)
	return v, nil
}

func LoadConfigKeyword_VERSIPOINTERWORK(filename string) (*map[string]map[string]string, error) {
	var v map[string]map[string]string
	configFile, err := os.Open(filename)
	if err != nil {
		models.ErrorLogger.Println(err)
		return nil, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&v)
	// models.InfoLogger.Println(len(v["HITACHI_KEYWORD"]))
	return &v, nil
}

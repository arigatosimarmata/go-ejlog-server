package utils

import (
	"ejol/ejlog-server/controller"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// func init() {
func InitUtils() {
	fmt.Println("access here first")
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error load file env : %s", err)
		// ErrorLogger.Fatal("Error load file env : %s", err)
	}
	filename := os.Getenv("EJOL_DIRECTORY_LOG") + "ejlog-server-" + time.Now().Format("20060102") + ".log"
	// filename := "ejlog-server-" + time.Now().Format("20060102") + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	controller.InfoLogger = log.New(file, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	controller.WarningLogger = log.New(file, "WARNING", log.Ldate|log.Ltime|log.Lshortfile)
	controller.ErrorLogger = log.New(file, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)

	// _, err := LoadConfigKeyword(os.Getenv("CONFIG_KEYWORD_PATH"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = k.LoadConfigKeyword(os.Getenv("CONFIG_KEYWORD_PATH"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	key, err2 := controller.LoadConfigKeyword(os.Getenv("CONFIG_KEYWORD_PATH"))
	if err2 != nil {
		log.Fatal(err2)
	}

	controller.KeywordEjol = &key
	// fmt.Println(*KeywordEjol)
	//HOW TO CALL USE THIS :
	// (*KeywordEjol)["HITACHI_KEYWORD"]

	// fmt.Println("==========================")
	// for key, value := range *KeywordEjol {
	// 	fmt.Println("Hasilnya")
	// 	fmt.Println(key, value)
	// }

}

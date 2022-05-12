package utils

import (
	"ejol/ejlog-server/controller"
	"ejol/ejlog-server/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// func init() {
func InitUtils() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error load file env : %s", err)
	}
	filename := os.Getenv("EJOL_DIRECTORY_LOG") + "ejlog-server-" + time.Now().Format("20060102") + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	models.InfoLogger = log.New(file, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	models.WarningLogger = log.New(file, "WARNING", log.Ldate|log.Ltime|log.Lshortfile)
	models.ErrorLogger = log.New(file, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)
	models.Logger = zerolog.New(file).With().Logger()

	key, err2 := controller.LoadConfigKeyword(os.Getenv("CONFIG_KEYWORD_PATH"))
	if err2 != nil {
		log.Fatal(err2)
	}

	unmapping := os.Getenv("UNMAPPING_FILE")
	models.KeywordEjol = &key
	models.Unmapping = &unmapping

	fmt.Println("Load Success.")
}

func InitUtilsDebug() {
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

	models.InfoLogger = log.New(file, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	models.WarningLogger = log.New(file, "WARNING", log.Ldate|log.Ltime|log.Lshortfile)
	models.ErrorLogger = log.New(file, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)
	models.Logger = zerolog.New(file).With().Logger()

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

	models.KeywordEjol = &key
	// fmt.Println(*KeywordEjol)
	//HOW TO CALL USE THIS :
	// (*KeywordEjol)["HITACHI_KEYWORD"]

	// fmt.Println("==========================")
	// for key, value := range *KeywordEjol {
	// 	fmt.Println("Hasilnya")
	// 	fmt.Println(key, value)
	// }

}

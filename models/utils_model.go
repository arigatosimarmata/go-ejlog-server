package models

import (
	"log"

	"github.com/rs/zerolog"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	Logger        zerolog.Logger
	KeywordEjol   *map[string]map[string]string
	Unmapping     *string
)

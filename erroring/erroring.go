package erroring

import (
	"log"
)

func QuitIfError(err error, message string) {
	if err != nil {
		defer panic(err)
		log.Printf("ERR | FATAL | %v | %v", message, err.Error())
	}
}

func LogIfError(err error, message string) {
	if err != nil {
		log.Printf("ERR | %v | %v", message, err.Error())
	}
}

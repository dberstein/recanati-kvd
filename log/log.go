package log

import (
	"log"
	"os"
)

func Print(v ...any) {
	log.Println(v...)
}

func Fatal(v ...any) {
	Print(v...)
	os.Exit(1)
}

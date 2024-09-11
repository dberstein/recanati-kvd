package log

import (
	"log"
	"os"
)

func Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func Print(v ...any) {
	log.Println(v...)
}

func Fatal(v ...any) {
	Print(v...)
	os.Exit(1)
}

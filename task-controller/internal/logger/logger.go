package logger

import (
	"log"
)

// Info логирует информационные сообщения.
func Info(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// Error логирует ошибки.
func Error(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

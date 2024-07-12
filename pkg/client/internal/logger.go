package internal

import (
	"log"
	"strings"
)

// LogInfo creates an info log.
func LogInfo(msg string, fields ...string) {
	if len(fields) > 0 {
		log.Printf("[info] %s: %s\n", msg, strings.Join(fields, ", "))
	} else {
		log.Printf("[info] %s\n", msg)
	}
}

// LogError creates an error log.
func LogError(msg string, e error, fields ...string) {
	if len(fields) > 0 {
		log.Printf("[error] %s: %s: %v\n", msg, strings.Join(fields, ", "), e)
	} else {
		log.Printf("[error] %s: %v\n", msg, e)
	}
}

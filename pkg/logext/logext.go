package logext

import "log"

func LogSuccess(err error, message string, args ...interface{}) {
	if err == nil {
		log.Printf("[SUCCESS] "+message+"\n", args...)
	}
}

func LogError(err error, prefix string, args ...interface{}) {
	if err != nil {
		log.Printf("[ERROR] "+prefix+": "+err.Error()+"\n", args...)
	}
}

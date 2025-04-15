package logext

import "log"

func LogIfSuccess(err error, message string, args ...interface{}) {
	if err == nil {
		LogSuccess(message, args...)
	}
}

func LogIfError(err error, prefix string, args ...interface{}) {
	if err != nil {
		LogError(prefix+": "+err.Error(), args...)
	}
}

func LogSuccess(message string, args ...interface{}) {
	log.Printf("[SUCCESS] "+message+"\n", args...)
}

func LogError(message string, args ...interface{}) {
	log.Printf("[ERROR] "+message+"\n", args...)
}

package logext

import "log"

func MaybeSuccess(err error, message string, args ...interface{}) {
	if err == nil {
		Success(message, args...)
	}
}

func MaybeWarning(err error, prefix string, args ...interface{}) {
	if err != nil {
		Warning(prefix+": "+err.Error(), args...)
	}
}

func MaybeError(err error, prefix string, args ...interface{}) {
	if err != nil {
		Error(prefix+": "+err.Error(), args...)
	}
}

func Info(message string, args ...interface{}) {
	log.Printf("[INFO]    "+message+"\n", args...)
}

func Success(message string, args ...interface{}) {
	log.Printf("[SUCCESS] "+message+"\n", args...)
}

func Warning(message string, args ...interface{}) {
	log.Printf("[WARNING] "+message+"\n", args...)
}

func Error(message string, args ...interface{}) {
	log.Printf("[ERROR]   "+message+"\n", args...)
}

func Request(url string) {
	log.Println("[REQUEST] (unauthorized) " + url)
}

func AuthorizedRequest(url string) {
	log.Println("[REQUEST] (authorized) " + url)
}

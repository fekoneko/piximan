package logext

import (
	"os"

	"github.com/fekoneko/piximan/internal/termext"
)

// TODO: multiline logs
// TODO: make a dictionary with log messages
// TODO: group logs if there are multiple of the same type

func Info(message string, args ...any) {
	log(infoPrefix+message, args...)
}

func Success(message string, args ...any) {
	log(successPrefix+message, args...)
}

func Warning(message string, args ...any) {
	log(warningPrefix+message, args...)
}

func Error(message string, args ...any) {
	log(errorPrefix+message, args...)
}

func Fatal(message string, args ...any) {
	log(errorPrefix+message, args...)
	termext.RestoreInputEcho()
	DisableProgress()
	os.Exit(1)
}

func MaybeSuccess(err error, message string, args ...any) {
	if err == nil {
		Success(message, args...)
	}
}

func MaybeWarning(err error, prefix string, args ...any) {
	if err != nil {
		Warning(prefix+": "+err.Error(), args...)
	}
}

func MaybeError(err error, prefix string, args ...any) {
	if err != nil {
		Error(prefix+": "+err.Error(), args...)
	}
}

func MaybeFatal(err error, prefix string, args ...any) {
	if err != nil {
		Fatal(prefix+": "+err.Error(), args...)
	}
}

func MaybeWarnings(errs []error, prefix string, args ...any) {
	for _, err := range errs {
		MaybeWarning(err, prefix, args...)
	}
}

func MaybeErrors(errs []error, prefix string, args ...any) {
	for _, err := range errs {
		MaybeError(err, prefix, args...)
	}
}

func Request(url string) (func(), func(int, int)) {
	removeBar, updateBar := registerRequest(url, false)
	log(requestPrefix + url)
	return removeBar, updateBar
}

func AuthorizedRequest(url string) (func(), func(int, int)) {
	removeBar, updateBar := registerRequest(url, true)
	log(authRequestPrefix + url)
	return removeBar, updateBar
}

func EnableProgress() {
	mutex.Lock()
	statsShown = true
	mutex.Unlock()
	refreshStats()
}

func DisableProgress() {
	mutex.Lock()
	statsShown = false
	mutex.Unlock()
	refreshStats()
}

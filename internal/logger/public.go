package logger

import (
	"os"

	"github.com/fekoneko/piximan/internal/termext"
)

// TODO: multiline logs
// TODO: make a dictionary with log messages
// TODO: group logs if there are multiple of the same type

func (l *Logger) Info(message string, args ...any) {
	l.log(infoPrefix+message, args...)
}

func (l *Logger) Success(message string, args ...any) {
	l.log(successPrefix+message, args...)
}

func (l *Logger) Warning(message string, args ...any) {
	l.log(warningPrefix+message, args...)
}

func (l *Logger) Error(message string, args ...any) {
	l.log(errorPrefix+message, args...)
}

func (l *Logger) Fatal(message string, args ...any) {
	l.log(errorPrefix+message, args...)
	termext.RestoreInputEcho()
	l.DisableProgress()
	os.Exit(1)
}

func (l *Logger) MaybeSuccess(err error, message string, args ...any) {
	if err == nil {
		l.Success(message, args...)
	}
}

func (l *Logger) MaybeWarning(err error, prefix string, args ...any) {
	if err != nil {
		l.Warning(prefix+": "+err.Error(), args...)
	}
}

func (l *Logger) MaybeError(err error, prefix string, args ...any) {
	if err != nil {
		l.Error(prefix+": "+err.Error(), args...)
	}
}

func (l *Logger) MaybeFatal(err error, prefix string, args ...any) {
	if err != nil {
		l.Fatal(prefix+": "+err.Error(), args...)
	}
}

func (l *Logger) MaybeWarnings(errs []error, prefix string, args ...any) {
	for _, err := range errs {
		l.MaybeWarning(err, prefix, args...)
	}
}

func (l *Logger) MaybeErrors(errs []error, prefix string, args ...any) {
	for _, err := range errs {
		l.MaybeError(err, prefix, args...)
	}
}

func (l *Logger) Request(url string) (removeBar func(), updateBar func(int, int)) {
	removeBar, updateBar = l.registerRequest(url, false)
	l.log(requestPrefix + url)
	return removeBar, updateBar
}

func (l *Logger) AuthorizedRequest(url string) (removeBar func(), updateBar func(int, int)) {
	removeBar, updateBar = l.registerRequest(url, true)
	l.log(authRequestPrefix + url)
	return removeBar, updateBar
}

func (l *Logger) EnableProgress() {
	l.mutex.Lock()
	l.statsShown = true
	l.mutex.Unlock()
	l.refreshStats()
}

func (l *Logger) DisableProgress() {
	l.mutex.Lock()
	l.statsShown = false
	l.mutex.Unlock()
	l.refreshStats()
}

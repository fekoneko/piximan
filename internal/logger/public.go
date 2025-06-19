package logger

import (
	"fmt"
	"os"
	"strings"

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
	l.mutex.Lock()
	l.numWarnings++
	l.mutex.Unlock()
	l.log(warningPrefix+message, args...)
}

func (l *Logger) Error(message string, args ...any) {
	l.mutex.Lock()
	l.numErrors++
	l.mutex.Unlock()
	l.log(errorPrefix+message, args...)
}

func (l *Logger) Fatal(message string, args ...any) {
	l.mutex.Lock()
	l.numErrors++
	l.mutex.Unlock()
	l.log(errorPrefix+message, args...)
	termext.RestoreInputEcho()
	l.HideProgress()
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

func (l *Logger) AddSuccessfulWork() {
	l.mutex.Lock()
	l.numSuccessfulWorks++
	l.mutex.Unlock()
}

func (l *Logger) AddFailedWork(id uint64) {
	l.mutex.Lock()
	l.failedWorkIds = append(l.failedWorkIds, id)
	l.mutex.Unlock()
}

func (l *Logger) AddSuccessfulCrawl() {
	l.mutex.Lock()
	l.numSuccessfulCrawls++
	l.mutex.Unlock()
}

func (l *Logger) AddFailedCrawl() {
	l.mutex.Lock()
	l.numFailedCrawls++
	l.mutex.Unlock()
}

func (l *Logger) ShowProgress() {
	l.mutex.Lock()
	l.progressShown = true
	l.mutex.Unlock()
	l.refreshProgress()
}

func (l *Logger) HideProgress() {
	l.mutex.Lock()
	l.progressShown = false
	l.mutex.Unlock()
	l.refreshProgress()
}

func (l *Logger) Stats() {
	l.mutex.Lock()

	builder := strings.Builder{}
	builder.WriteString("\ndownloader stats:\n")
	builder.WriteString(fmt.Sprintf("- unauthorized requests: %-6v\n", l.numRequests-l.numAuthorizedRequests))
	builder.WriteString(fmt.Sprintf("- authorized requests: %-6v\n", l.numAuthorizedRequests))
	builder.WriteString(fmt.Sprintf("- total requests: %-6v\n", l.numRequests))
	s := fmt.Sprintf("- warnings: %-6v\n", l.numWarnings)
	if l.numWarnings > 0 {
		s = yellow(s)
	}
	builder.WriteString(s)
	s = fmt.Sprintf("- errors: %-6v\n", l.numErrors)
	if l.numErrors > 0 {
		s = red(s)
	}
	builder.WriteString(s)
	builder.WriteString(fmt.Sprintf("- works downloaded: %-6v\n", l.numSuccessfulWorks))
	s = fmt.Sprintf("- failed works: %-6v\n", len(l.failedWorkIds))
	if len(l.failedWorkIds) > 0 {
		s = red(s)
	}
	builder.WriteString(s)
	// TODO: show the failed work ids
	builder.WriteString(fmt.Sprintf("- successful crawl tasks: %-6v\n", l.numSuccessfulCrawls))
	s = fmt.Sprintf("- failed crawl tasks: %-6v\n", l.numFailedCrawls)
	if l.numFailedCrawls > 0 {
		s = red(s)
	}
	builder.WriteString(s)

	l.mutex.Unlock()
	l.printWithProgress("%v", builder.String())
}

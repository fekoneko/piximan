package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/fekoneko/piximan/internal/termext"
)

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

func (l *Logger) MaybeInfo(err error, message string, args ...any) {
	if err == nil {
		l.Info(message, args...)
	}
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

func (l *Logger) Request(url string) (RemoveBarFunc, UpdateBarFunc) {
	removeBar, updateBar := l.registerRequest(url, false)
	l.log("%v%v", requestPrefix, url)
	return removeBar, updateBar
}

func (l *Logger) AuthorizedRequest(url string) (RemoveBarFunc, UpdateBarFunc) {
	removeBar, updateBar := l.registerRequest(url, true)
	l.log("%v%v", authRequestPrefix, url)
	return removeBar, updateBar
}

func (l *Logger) ExpectWorks(count int) {
	l.mutex.Lock()
	l.numExpectedWorks += count
	l.mutex.Unlock()
}

func (l *Logger) AddSuccessfulWork() {
	l.mutex.Lock()
	l.numSuccessfulWorks++
	l.mutex.Unlock()
}

func (l *Logger) AddSkippedWork() {
	l.mutex.Lock()
	l.numSkippedWorks++
	l.mutex.Unlock()
}

func (l *Logger) AddFailedWork(id uint64) { // TODO: pass work's queue.ItemKind as well
	l.mutex.Lock()
	l.failedWorkIds = append(l.failedWorkIds, id)
	l.mutex.Unlock()
}

func (l *Logger) ExpectCrawls(count int) {
	l.mutex.Lock()
	l.numExpectedCrawls += count
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

func (l *Logger) AddSkippedCrawl() {
	l.mutex.Lock()
	l.numSkippedCrawls++
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
	builder := strings.Builder{}
	builder.WriteByte('\n')
	l.mutex.Lock()

	numTotalCrawls := l.numExpectedCrawls - l.numSkippedCrawls
	builder.WriteString(fmt.Sprintf("- tasks crawled: %v / %v", l.numSuccessfulCrawls, numTotalCrawls))
	if l.numSkippedCrawls > 0 {
		builder.WriteString(fmt.Sprintf(" + %v skipped", l.numSkippedCrawls))
	}

	numTotalWorks := l.numExpectedWorks - l.numSkippedWorks
	builder.WriteString(fmt.Sprintf("\n- works downloaded: %v / %v", l.numSuccessfulWorks, numTotalWorks))
	if l.numSkippedWorks > 0 {
		builder.WriteString(fmt.Sprintf(" + %v skipped", l.numSkippedWorks))
	}

	builder.WriteString(fmt.Sprintf("\n- unauthorized requests: %v\n", l.numRequests-l.numAuthorizedRequests))
	builder.WriteString(fmt.Sprintf("- authorized requests: %v\n\n", l.numAuthorizedRequests))

	s := fmt.Sprintf("- warnings: %v\n", l.numWarnings)
	if l.numWarnings > 0 {
		s = yellow(s)
	}
	builder.WriteString(s)
	s = fmt.Sprintf("- errors: %v\n", l.numErrors)
	if l.numErrors > 0 {
		s = red(s)
	}
	builder.WriteString(s)
	s = fmt.Sprintf("- failed crawl tasks: %v\n", l.numFailedCrawls)
	if l.numFailedCrawls > 0 {
		s = red(s)
	}
	builder.WriteString(s)
	s = fmt.Sprintf("- failed works: %v", len(l.failedWorkIds))
	if l.numFailedCrawls > 0 {
		s += " (?)"
	}
	if len(l.failedWorkIds) > 0 {
		s = red(s)
	} else if l.numFailedCrawls > 0 {
		s = yellow(s)
	}
	builder.WriteString(s)

	if len(l.failedWorkIds) > 0 {
		idsBuilder := strings.Builder{}
		for i, id := range l.failedWorkIds {
			if i%6 == 0 {
				idsBuilder.WriteString("\n  | ")
			}
			idsBuilder.WriteString(fmt.Sprintf("%-11v", id))
		}
		builder.WriteString(red(idsBuilder.String()))
	}
	builder.WriteByte('\n')

	l.mutex.Unlock()
	l.printWithProgress(builder.String())
}

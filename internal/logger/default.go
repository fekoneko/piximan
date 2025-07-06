package logger

import "os"

var DefaultLogger = New(os.Stdout)
var l = DefaultLogger

func Info(message string, args ...any)                       { l.Info(message, args...) }
func Success(message string, args ...any)                    { l.Success(message, args...) }
func Warning(message string, args ...any)                    { l.Warning(message, args...) }
func Error(message string, args ...any)                      { l.Error(message, args...) }
func Fatal(message string, args ...any)                      { l.Fatal(message, args...) }
func MaybeSuccess(err error, message string, args ...any)    { l.MaybeSuccess(err, message, args...) }
func MaybeWarning(err error, prefix string, args ...any)     { l.MaybeWarning(err, prefix, args...) }
func MaybeError(err error, prefix string, args ...any)       { l.MaybeError(err, prefix, args...) }
func MaybeFatal(err error, prefix string, args ...any)       { l.MaybeFatal(err, prefix, args...) }
func MaybeWarnings(errs []error, prefix string, args ...any) { l.MaybeWarnings(errs, prefix, args...) }
func MaybeErrors(errs []error, prefix string, args ...any)   { l.MaybeErrors(errs, prefix, args...) }
func ExpectWorks(count int)                                  { l.ExpectWorks(count) }
func AddSuccessfulWork()                                     { l.AddSuccessfulWork() }
func AddFailedWork(id uint64)                                { l.AddFailedWork(id) }
func ExpectCrawls(count int)                                 { l.ExpectCrawls(count) }
func AddSuccessfulCrawl()                                    { l.AddSuccessfulCrawl() }
func AddFailedCrawl()                                        { l.AddFailedCrawl() }
func EnableProgress()                                        { l.ShowProgress() }
func DisableProgress()                                       { l.HideProgress() }
func Stats()                                                 { l.Stats() }

func Request(url string) (removeBar func(), updateBar func(int, int)) {
	return l.Request(url)
}
func AuthorizedRequest(url string) (removeBar func(), updateBar func(int, int)) {
	return l.AuthorizedRequest(url)
}

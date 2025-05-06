package logext

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// TODO: separate files

var cyan = color.New(color.FgHiCyan, color.Bold).SprintFunc()
var green = color.New(color.FgHiGreen, color.Bold).SprintFunc()
var yellow = color.New(color.FgHiYellow, color.Bold).SprintFunc()
var red = color.New(color.FgHiRed, color.Bold).SprintFunc()
var magenta = color.New(color.FgHiMagenta, color.Bold).SprintFunc()
var white = color.New(color.FgHiWhite, color.Bold).SprintFunc()
var gray = color.New(color.FgHiBlack).SprintFunc()

var infoPrefix = cyan("[INFO]") + "    "
var successPrefix = green("[SUCCESS]") + " "
var warningPrefix = yellow("[WARNING]") + " "
var errorPrefix = red("[ERROR]") + "   "
var requestPrefix = magenta("[REQUEST]") + " " + white("(unauthorized)") + " "
var authRequestPrefix = magenta("[REQUEST]") + " " + red("(authorized)") + " "

var mutex = sync.Mutex{}
var requests = map[uint]*request{}
var numRequestsDisplayed = int(0)
var requestsTotal = uint(0)

type request struct {
	url     string
	current int
	total   int
}

func (r *request) String() string {
	percent := 0
	if r.total > 0 && r.current > 0 {
		percent = int(float64(r.current) / float64(r.total) * 100)
	}
	return fmt.Sprintf("%+3v%% | %v", percent, r.url)
}

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
	settle, updateProgress := handleRequest(url)
	log(requestPrefix + url)
	return settle, updateProgress
}

func AuthorizedRequest(url string) (func(), func(int, int)) {
	settle, updateProgress := handleRequest(url)
	log(authRequestPrefix + url)
	return settle, updateProgress
}

func log(message string, args ...any) {
	timePrefix := gray(time.Now().Format(time.DateTime)) + " "
	print(timePrefix+message+"\n", args...)
}

func print(s string, args ...any) {
	mutex.Lock()
	defer mutex.Unlock()

	if numRequestsDisplayed > 0 {
		for range numRequestsDisplayed + 1 {
			fmt.Print("\033[2K\033[A\033[2K\r")
		}
	}

	fmt.Fprintf(color.Output, s, args...)

	if len(requests) > 0 {
		fmt.Fprintln(color.Output)
	}
	for _, request := range requests {
		fmt.Fprintf(color.Output, gray("%v\n"), request)
	}
	numRequestsDisplayed = len(requests)
}

func handleRequest(url string) (func(), func(int, int)) {
	mutex.Lock()
	requestsTotal++
	requestIndex := requestsTotal
	requests[requestIndex] = &request{url, 0, -1}
	mutex.Unlock()

	settle := func() {
		mutex.Lock()
		delete(requests, requestIndex)
		mutex.Unlock()
	}

	updateProgress := func(current int, total int) {
		mutex.Lock()
		requests[requestIndex].current = current
		requests[requestIndex].total = total
		mutex.Unlock()
		print("")
	}

	return settle, updateProgress
}

package logext

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// TODO: separate files

const REQUEST_URL_LENGTH = 36
const BAR_LENGTH = 36

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
var requests = map[int]*request{}
var numRequests = int(0)
var requestsTotal = int(0)

type request struct {
	url     string
	current int
	total   int
}

func (r *request) String() string {
	var url string
	domainStart := strings.Index(r.url, "://") + 3
	if domainStart == 2 {
		domainStart = 0
	}
	domainEnd := strings.Index(r.url[domainStart:], "/")
	if domainEnd == -1 {
		url = r.url[domainStart:]
	} else {
		domainEnd += domainStart
		domain := r.url[domainStart:domainEnd]
		suffixStart := len(r.url) - (REQUEST_URL_LENGTH - 4 - len(domain))
		if suffixStart-domainEnd <= 4 {
			url = r.url[domainStart:]
		} else {
			url = domain + "/..." + r.url[suffixStart:]
		}
	}

	bar := barString(r.current, r.total)
	return fmt.Sprintf(gray("%-*v "), REQUEST_URL_LENGTH, url) + bar
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
	removeBar, updateBar := handleRequest(url)
	log(requestPrefix + url)
	return removeBar, updateBar
}

func AuthorizedRequest(url string) (func(), func(int, int)) {
	removeBar, updateBar := handleRequest(url)
	log(authRequestPrefix + url)
	return removeBar, updateBar
}

func log(message string, args ...any) {
	timePrefix := gray(time.Now().Format(time.DateTime)) + " "
	printWithRequests(timePrefix+message+"\n", args...)
}

func printWithRequests(s string, args ...any) {
	mutex.Lock()
	defer mutex.Unlock()

	builder := strings.Builder{}

	if numRequests > 0 {
		for range numRequests + 1 {
			builder.WriteString("\033[2K\033[A\033[2K\r")
		}
	}
	builder.WriteString(fmt.Sprintf(s, args...))
	if len(requests) > 0 {
		builder.WriteByte('\n')
	}
	indeces := make([]int, 0, len(requests))
	for index := range requests {
		indeces = append(indeces, index)
	}
	sort.Ints(indeces)
	for _, index := range indeces {
		builder.WriteString(requests[index].String())
		builder.WriteByte('\n')
	}

	fmt.Fprint(color.Output, builder.String())
	numRequests = len(requests)
}

func handleRequest(url string) (func(), func(int, int)) {
	mutex.Lock()
	requestsTotal++
	requestIndex := requestsTotal
	requests[requestIndex] = &request{url, 0, -1}
	mutex.Unlock()

	removeBar := func() {
		mutex.Lock()
		delete(requests, requestIndex)
		mutex.Unlock()
	}

	updateBar := func(current int, total int) {
		mutex.Lock()
		requests[requestIndex].current = current
		requests[requestIndex].total = total
		mutex.Unlock()
		printWithRequests("")
	}

	return removeBar, updateBar
}

func barString(current int, total int) string {
	fraction := float64(0)
	if total > 0 && current > 0 {
		fraction = float64(current) / float64(total)
	}
	percent := int(math.Round(fraction * 100))
	chars := int(math.Round(fraction * float64(BAR_LENGTH)))
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf(gray("%3v%% "), percent))

	for i := 0; i < BAR_LENGTH; i++ {
		if i < chars {
			builder.WriteString(white("━"))
		} else if i == chars && i != 0 {
			builder.WriteString(gray("╶"))
		} else {
			builder.WriteString(gray("─"))
		}
	}

	return builder.String()
}

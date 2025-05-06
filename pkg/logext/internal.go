package logext

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

const NUM_REQUEST_SLOTS = 6
const REQUEST_URL_LENGTH = 36
const BAR_LENGTH = 36

var cyan = color.New(color.FgHiCyan, color.Bold).SprintFunc()
var green = color.New(color.FgHiGreen, color.Bold).SprintFunc()
var yellow = color.New(color.FgHiYellow, color.Bold).SprintFunc()
var red = color.New(color.FgHiRed, color.Bold).SprintFunc()
var magenta = color.New(color.FgHiMagenta, color.Bold).SprintFunc()
var white = color.New(color.FgHiWhite, color.Bold).SprintFunc()
var gray = color.New(color.FgHiBlack, color.Bold).SprintFunc()
var subtleGray = color.New(color.FgHiBlack).SprintFunc()

var infoPrefix = cyan("[INFO]") + "    "
var successPrefix = green("[SUCCESS]") + " "
var warningPrefix = yellow("[WARNING]") + " "
var errorPrefix = red("[ERROR]") + "   "
var requestPrefix = magenta("[REQUEST]") + " " + white("(unauthorized)") + " "
var authRequestPrefix = magenta("[REQUEST]") + " " + red("(authorized)") + " "

var mutex = sync.Mutex{}
var requests = map[int]*request{}
var requestSlots = make([]int, NUM_REQUEST_SLOTS)
var requestSlotsShown = false
var prevRequestSlotsShown = false
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

func log(message string, args ...any) {
	timePrefix := subtleGray(time.Now().Format(time.DateTime)) + " "
	printWithRequestSlots(timePrefix+message+"\n", args...)
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
		printWithRequestSlots("")
	}

	updateBar := func(current int, total int) {
		mutex.Lock()
		requests[requestIndex].current = current
		requests[requestIndex].total = total
		mutex.Unlock()
		printWithRequestSlots("")
	}

	return removeBar, updateBar
}

func printWithRequestSlots(s string, args ...any) {
	builder := strings.Builder{}

	mutex.Lock()
	defer mutex.Unlock()

	if prevRequestSlotsShown {
		for range NUM_REQUEST_SLOTS + 1 {
			builder.WriteString("\033[2K\033[A\033[2K\r")
		}
	}

	if !requestSlotsShown {
		prevRequestSlotsShown = false
		builder.WriteString(fmt.Sprintf(s, args...))
		fmt.Fprint(color.Output, builder.String())
		return
	}

	builder.WriteString(fmt.Sprintf(s, args...))

	builder.WriteByte('\n')
	hasNewRequests := true
	for i, index := range requestSlots {
		if request, ok := requests[index]; ok {
			addSlotToBuilder(&builder, request)
			continue
		}
		requestSlots[i] = 0
		var request *request
		if hasNewRequests {
			for index := range requests {
				if !slices.Contains(requestSlots, index) {
					requestSlots[i] = index
					request = requests[index]
					hasNewRequests = false
					break
				}
			}
		}
		addSlotToBuilder(&builder, request)
	}

	prevRequestSlotsShown = true
	fmt.Fprint(color.Output, builder.String())
}

func addSlotToBuilder(builder *strings.Builder, request *request) {
	if request == nil {
		for range REQUEST_URL_LENGTH + BAR_LENGTH + 6 {
			builder.WriteString(subtleGray("╶"))
		}
		builder.WriteByte('\n')
	} else {
		builder.WriteString(request.String())
		builder.WriteByte('\n')
	}
}

func barString(current int, total int) string {
	fraction := float64(0)
	if total > 0 && current > 0 {
		fraction = float64(current) / float64(total)
	}
	percent := int(math.Round(fraction * 100))
	chars := int(math.Round(fraction * float64(BAR_LENGTH)))
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf(subtleGray("%3v%% "), percent))

	for i := 0; i < BAR_LENGTH; i++ {
		if i < chars {
			builder.WriteString(white("━"))
		} else if i == chars && i != 0 {
			builder.WriteString(subtleGray("╶"))
		} else {
			builder.WriteString(subtleGray("─"))
		}
	}

	return builder.String()
}

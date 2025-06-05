package logext

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/fatih/color"
)

func log(message string, args ...any) {
	timePrefix := subtleGray(time.Now().Format(time.DateTime)) + " "
	printWithStats(timePrefix+message+"\n", args...)
}

// track request internally and return handlers to update its stats
func registerRequest(url string, authorized bool) (func(), func(int, int)) {
	mutex.Lock()
	numRequests++
	mapIndex := numRequests
	progressMap[mapIndex] = &progress{url, 0, -1}
	if authorized {
		numAuthorizedRequests++
	}
	mutex.Unlock()

	removeBar := func() {
		mutex.Lock()
		delete(progressMap, mapIndex)
		mutex.Unlock()
		refreshStats()
	}

	updateBar := func(current int, total int) {
		mutex.Lock()
		progressMap[mapIndex].current = current
		progressMap[mapIndex].total = total
		mutex.Unlock()
		refreshStats()
	}

	return removeBar, updateBar
}

func printWithStats(s string, args ...any) {
	builder := strings.Builder{}

	mutex.Lock()
	defer mutex.Unlock()

	if prevStatsShown {
		for range NUM_SLOTS + 2 {
			builder.WriteString("\033[2K\033[A\033[2K\r")
		}
	}

	if !statsShown {
		prevStatsShown = false
		builder.WriteString(fmt.Sprintf(s, args...))
		fmt.Fprint(color.Output, builder.String())
		return
	}

	builder.WriteString(fmt.Sprintf(s, args...))
	builder.WriteByte('\n')

	hasNewRequests := true
	for i, index := range slots {
		if request, ok := progressMap[index]; ok {
			addSlot(&builder, request)
			continue
		}
		slots[i] = 0
		var request *progress
		if hasNewRequests {
			for index := range progressMap {
				if !slices.Contains(slots, index) {
					slots[i] = index
					request = progressMap[index]
					hasNewRequests = false
					break
				}
			}
		}
		addSlot(&builder, request)
	}

	builder.WriteString(fmt.Sprintf(
		"pending: %-6v unauthorized: %-6v authorized: %-6v total: %v\n",
		len(progressMap), numRequests-numAuthorizedRequests, numAuthorizedRequests, numRequests),
	)

	prevStatsShown = true
	fmt.Fprint(color.Output, builder.String())
}

func addSlot(builder *strings.Builder, request *progress) {
	if request == nil {
		for range URL_LENGTH + BAR_LENGTH + 6 {
			builder.WriteString(subtleGray("╶"))
		}
		builder.WriteByte('\n')
	} else {
		builder.WriteString(request.String())
		builder.WriteByte('\n')
	}
}

func refreshStats() {
	printWithStats("")
}

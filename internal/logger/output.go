package logger

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/fatih/color"
)

func (l *Logger) log(message string, args ...any) {
	timePrefix := subtleGray(time.Now().Format(time.DateTime)) + " "
	l.printWithStats(timePrefix+message+"\n", args...)
}

// track request internally and return handlers to update its state
func (l *Logger) registerRequest(url string, authorized bool) (func(), func(int, int)) {
	l.mutex.Lock()
	l.numRequests++
	mapIndex := l.numRequests
	l.progressMap[mapIndex] = &progress{url, 0, -1}
	if authorized {
		l.numAuthorizedRequests++
	}
	l.mutex.Unlock()

	removeBar := func() {
		l.mutex.Lock()
		delete(l.progressMap, mapIndex)
		l.mutex.Unlock()
		l.refreshStats()
	}

	updateBar := func(current int, total int) {
		l.mutex.Lock()
		l.progressMap[mapIndex].current = current
		l.progressMap[mapIndex].total = total
		l.mutex.Unlock()
		l.refreshStats()
	}

	return removeBar, updateBar
}

func (l *Logger) printWithStats(s string, args ...any) {
	builder := strings.Builder{}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.prevStatsShown {
		for range NUM_SLOTS + 2 {
			builder.WriteString("\033[2K\033[A\033[2K\r")
		}
	}

	if !l.statsShown {
		l.prevStatsShown = false
		builder.WriteString(fmt.Sprintf(s, args...))
		fmt.Fprint(color.Output, builder.String())
		return
	}

	builder.WriteString(fmt.Sprintf(s, args...))
	builder.WriteByte('\n')

	hasNewRequests := true
	for i, index := range l.slots {
		if request, ok := l.progressMap[index]; ok {
			addSlot(&builder, request)
			continue
		}
		l.slots[i] = 0
		var request *progress
		if hasNewRequests {
			for index := range l.progressMap {
				if !slices.Contains(l.slots, index) {
					l.slots[i] = index
					request = l.progressMap[index]
					hasNewRequests = false
					break
				}
			}
		}
		addSlot(&builder, request)
	}

	builder.WriteString(fmt.Sprintf(
		"pending: %-6v unauthorized: %-6v authorized: %-6v total: %v\n",
		len(l.progressMap), l.numRequests-l.numAuthorizedRequests, l.numAuthorizedRequests, l.numRequests),
	)

	l.prevStatsShown = true

	fmt.Fprint(*l.writer, builder.String())
}

func (l *Logger) refreshStats() {
	l.printWithStats("")
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

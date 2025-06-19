package logger

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/fatih/color"
)

func (l *Logger) log(message string, args ...any) {
	timePrefix := subtleGray(time.Now().Format(time.DateTime)) + " "
	l.printWithProgress(timePrefix+message+"\n", args...)
}

// track request internally and return handlers to update its state
func (l *Logger) registerRequest(url string, authorized bool) (removeBar func(), updateBar func(int, int)) {
	l.mutex.Lock()
	l.numRequests++
	mapIndex := l.numRequests
	l.progressMap[mapIndex] = &progress{url, 0, -1}
	if authorized {
		l.numAuthorizedRequests++
	}
	l.mutex.Unlock()

	removeBar = func() {
		l.mutex.Lock()
		delete(l.progressMap, mapIndex)
		l.mutex.Unlock()
		l.refreshProgress()
	}

	updateBar = func(current int, total int) {
		l.mutex.Lock()
		l.progressMap[mapIndex].current = current
		l.progressMap[mapIndex].total = total
		l.mutex.Unlock()
		l.refreshProgress()
	}

	return removeBar, updateBar
}

func (l *Logger) printWithProgress(s string, args ...any) {
	builder := strings.Builder{}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.prevProgressShown {
		eraseProgress(&builder)
	}
	if !l.progressShown {
		l.prevProgressShown = false
		builder.WriteString(fmt.Sprintf(s, args...))
		fmt.Fprint(color.Output, builder.String())
	} else {
		builder.WriteString(fmt.Sprintf(s, args...))
		builder.WriteByte('\n')

		l.addSlots(&builder)
		l.addStats(&builder)
		l.prevProgressShown = true
		fmt.Fprint(*l.writer, builder.String())
	}
}

func (l *Logger) refreshProgress() {
	l.printWithProgress("")
}

func (l *Logger) addSlots(builder *strings.Builder) {
	hasNewRequests := true
	for i, index := range l.slots {
		if progress, ok := l.progressMap[index]; ok {
			addSlot(builder, progress)
			continue
		}
		l.slots[i] = 0
		var progress *progress
		if hasNewRequests {
			for index := range l.progressMap {
				if !slices.Contains(l.slots, index) {
					l.slots[i] = index
					progress = l.progressMap[index]
					hasNewRequests = false
					break
				}
			}
		}
		addSlot(builder, progress)
	}
}

func eraseProgress(builder *strings.Builder) {
	for range NUM_SLOTS + 4 {
		builder.WriteString("\033[2K\033[A\033[2K\r")
	}
}

func addSlot(builder *strings.Builder, progress *progress) {
	if progress == nil {
		for range URL_LENGTH + BAR_LENGTH + 6 {
			builder.WriteString(subtleGray("╶"))
		}
		builder.WriteByte('\n')
	} else {
		builder.WriteString(progress.String())
		builder.WriteByte('\n')
	}
}

func (l *Logger) addStats(builder *strings.Builder) {
	const captionsLength = 26
	const length = BAR_LENGTH + URL_LENGTH - captionsLength

	numSettledCrawls := l.numSuccessfulCrawls + l.numFailedCrawls
	s := fmt.Sprintf("crawling (%v / %v): ", numSettledCrawls, l.numExpectedCrawls)
	builder.WriteString(fmt.Sprintf(gray("%-*v "), captionsLength, s))
	bar := barString(numSettledCrawls, l.numExpectedCrawls, length)
	builder.WriteString(bar)
	builder.WriteByte('\n')

	numSettledWorks := l.numSuccessfulWorks + len(l.failedWorkIds)
	s = fmt.Sprintf("downloading (%v / %v): ", numSettledWorks, l.numExpectedWorks)
	builder.WriteString(fmt.Sprintf(gray("%-*v "), captionsLength, s))
	bar = barString(numSettledWorks, l.numExpectedWorks, length)
	builder.WriteString(bar)
	builder.WriteByte('\n')

	builder.WriteString(fmt.Sprintf(
		"pending: %-6v unauthorized: %-6v authorized: %-6v total: %v\n",
		len(l.progressMap), l.numRequests-l.numAuthorizedRequests, l.numAuthorizedRequests, l.numRequests,
	))
}

func barString(current int, total int, length int) string {
	fraction := float64(1)
	if total > 0 {
		fraction = float64(current) / float64(total)
	}
	percent := int(math.Round(fraction * 100))
	chars := int(math.Round(fraction * float64(length)))
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf(subtleGray("%3v%% "), percent))

	for i := 0; i < length; i++ {
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

type progress struct {
	url     string
	current int
	total   int
}

func (r *progress) String() string {
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
		suffixStart := len(r.url) - (URL_LENGTH - 4 - len(domain))
		if suffixStart-domainEnd <= 4 {
			url = r.url[domainStart:]
		} else {
			url = domain + "/..." + r.url[suffixStart:]
		}
	}

	bar := barString(r.current, r.total, BAR_LENGTH)
	return fmt.Sprintf(gray("%-*v "), URL_LENGTH, url) + bar
}

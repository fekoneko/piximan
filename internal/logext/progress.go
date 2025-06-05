package logext

import (
	"fmt"
	"math"
	"strings"
)

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

	bar := barString(r.current, r.total)
	return fmt.Sprintf(gray("%-*v "), URL_LENGTH, url) + bar
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

package logext

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var cyan = color.New(color.FgHiCyan, color.Bold).SprintFunc()
var green = color.New(color.FgHiGreen, color.Bold).SprintFunc()
var yellow = color.New(color.FgHiYellow, color.Bold).SprintFunc()
var red = color.New(color.FgHiRed, color.Bold).SprintFunc()
var magenta = color.New(color.FgHiMagenta, color.Bold).SprintFunc()
var white = color.New(color.FgHiWhite, color.Bold).SprintFunc()

var infoPrefix = cyan("[INFO]") + "    "
var successPrefix = green("[SUCCESS]") + " "
var warningPrefix = yellow("[WARNING]") + " "
var errorPrefix = red("[ERROR]") + "   "
var requestPrefix = magenta("[REQUEST]") + " " + white("(unauthorized)") + " "
var authRequestPrefix = magenta("[REQUEST]") + " " + red("(authorized)") + " "

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

func Info(message string, args ...any) {
	fmt.Fprintf(color.Output, timePrefix()+infoPrefix+message+"\n", args...)
}

func Success(message string, args ...any) {
	fmt.Fprintf(color.Output, timePrefix()+successPrefix+message+"\n", args...)
}

func Warning(message string, args ...any) {
	fmt.Fprintf(color.Output, timePrefix()+warningPrefix+message+"\n", args...)
}

func Error(message string, args ...any) {
	fmt.Fprintf(color.Output, timePrefix()+errorPrefix+message+"\n", args...)
}

func Request(url string) {
	fmt.Fprintln(color.Output, timePrefix()+requestPrefix+url)
}

func AuthorizedRequest(url string) {
	fmt.Fprintln(color.Output, timePrefix()+authRequestPrefix+url)
}

func timePrefix() string {
	return time.Now().Format(time.DateTime) + " "
}

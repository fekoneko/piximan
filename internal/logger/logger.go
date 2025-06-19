package logger

import (
	"io"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

const NUM_SLOTS = 6
const URL_LENGTH = 36
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

// Used to log the messages and display request statuses.
// Avoid using multiple loggers on the same output at the same time.
type Logger struct {
	mutex                 *sync.Mutex
	writer                *io.Writer
	progressMap           map[int]*progress
	slots                 []int
	progressShown         bool
	prevProgressShown     bool
	numRequests           int
	numAuthorizedRequests int
	numSuccessfulWorks    int
	failedWorkIds         []uint64
	numSuccessfulCrawls   int
	numFailedCrawls       int
	numWarnings           int
	numErrors             int
}

func New(output *os.File) *Logger {
	writer := colorable.NewColorable(output)

	return &Logger{
		mutex:         &sync.Mutex{},
		writer:        &writer,
		progressMap:   map[int]*progress{},
		slots:         make([]int, NUM_SLOTS),
		failedWorkIds: make([]uint64, 0),
	}
}

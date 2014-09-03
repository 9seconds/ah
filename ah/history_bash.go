package ah

// --- Imports

import (
	"strconv"
	"strings"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

// --- Vars

var timestampRegex = pcre.MustCompile(`^#'s*\d+$`, 0)

// --- Structs

type historyScannerBash struct {
	historyScanner
}

// --- Methods

func (hsb *historyScannerBash) GetCommands() ([]HistoryCommand, error) {
	commandTime := TimeNow

	for hsb.scanner.Scan() {
		text := strings.TrimSpace(hsb.scanner.Text())
		matcher := timestampRegex.MatcherString(text, 0)

		if matcher.MatchString(text, 0) {
			if parsed, err := strconv.Atoi(text[1:]); err == nil {
				commandTime = convertTimestamp(parsed)
			}
		} else {
			hsb.addCommand(text, commandTime)
		}
	}

	return hsb.finishScan()
}

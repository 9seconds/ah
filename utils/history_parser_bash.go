package utils

import (
	"time"
	"strings"
	"strconv"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

var timestampRegex = pcre.MustCompile(`^#'s*\d+$`, 0)

type historyScannerBash struct {
	historyScanner
}


func (hsb *historyScannerBash) GetCommands() ([]HistoryCommand, error) {
	commandTime := TimeNow

	for hsb.scanner.Scan() {
		text := strings.TrimSpace(hsb.scanner.Text())
		matcher := timestampRegex.MatcherString(text, 0)

		if matcher.MatchString(text, 0) {
			if parsed, err := strconv.Atoi(text[1:]); err == nil {
				commandTime = time.Unix(int64(parsed), 0)
			}
		} else {
			parsedCommand := HistoryCommand{text, commandTime}
			hsb.content = append(hsb.content, parsedCommand)
		}
	}

	if err := hsb.scanner.Err(); err != nil {
		return nil, err
	}

	return hsb.content, nil
}

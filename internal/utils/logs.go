package utils

import (
	"fmt"
	"os"
	"strings"
)

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func WrapYellow(msg ...string) string {
	return wrapLog("33", msg...)
}

func WrapGrey(msg ...string) string {
	return wrapLog("37", msg...)
}

func WrapRed(msg ...string) string {
	return wrapLog("91", msg...)
}

func wrapLog(color string, msg ...string) string {
	if isTerminal() {
		return fmt.Sprintf("\033[1;"+color+"m%v\033[0m", strings.Join(msg, ""))
	} else {
		return fmt.Sprintf(strings.Join(msg, ""))
	}
}

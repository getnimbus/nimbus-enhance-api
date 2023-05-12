package util

import (
	"regexp"
)

var removeSpecialCharsRegex = regexp.MustCompile(`[^a-zA-Z0-9()@:%_\+.~#?&//=\- ]+`)

func RemoveSpecialChars(str string) string {
	return removeSpecialCharsRegex.ReplaceAllString(str, "")
}

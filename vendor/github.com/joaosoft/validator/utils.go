package validator

import (
	"regexp"
	"strings"
)

var (
	replaces = map[*regexp.Regexp]string{
		regexp.MustCompile(`[\xC0-\xC6]`): "A",
		regexp.MustCompile(`[\xC0-\xC6]`): "A",
		regexp.MustCompile(`[\xE0-\xE6]`): "a",
		regexp.MustCompile(`[\xC8-\xCB]`): "E",
		regexp.MustCompile(`[\xE8-\xEB]`): "e",
		regexp.MustCompile(`[\xCC-\xCF]`): "I",
		regexp.MustCompile(`[\xEC-\xEF]`): "i",
		regexp.MustCompile(`[\xD2-\xD6]`): "O",
		regexp.MustCompile(`[\xF2-\xF6]`): "o",
		regexp.MustCompile(`[\xD9-\xDC]`): "U",
		regexp.MustCompile(`[\xF9-\xFC]`): "u",
		regexp.MustCompile(`[\xC7-\xE7]`): "c",
		regexp.MustCompile(`[\xD1]`):      "N",
		regexp.MustCompile(`[\xF1]`):      "n",
	}
	spacereg       = regexp.MustCompile(`\s+`)
	noncharreg     = regexp.MustCompile(`[^A-Za-z0-9-]`)
	minusrepeatreg = regexp.MustCompile(`\-{2,}`)
)

func convertToKey(str string, lower bool) string {
	for regex, replace := range replaces {
		str = regex.ReplaceAllString(str, replace)
	}

	if lower {
		str = strings.ToLower(str)
	}
	str = spacereg.ReplaceAllString(str, "-")
	str = noncharreg.ReplaceAllString(str, "")
	str = minusrepeatreg.ReplaceAllString(str, "-")

	return str
}

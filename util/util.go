package util

import (
	"regexp"
	"strconv"
	"strings"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func InttoString[T Number](elem T) string {
	return strconv.FormatInt(int64(elem), 10)
}

func IDtoString(elem int64) string {
	return strconv.FormatInt(elem, 10)
}
func JoinInt[T Number](elems *[]T, sep *string) string {
	// default value
	if *sep == "" {
		*sep = ", "
	}
	switch len(*elems) {
	case 0:
		return ""
	case 1:
		return InttoString((*elems)[0])
	}
	n := len(*sep) * (len(*elems) - 1)
	for i := 0; i < len(*elems); i++ {
		current := InttoString((*elems)[i])
		n += len(current)
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(InttoString((*elems)[0]))
	for index, num := range *elems {
		if index == 0 {
			continue
		}
		b.WriteString(*sep)
		b.WriteString(InttoString(num))
	}
	return b.String()
}

func AssignDefault(value *string, defaultValue string) {
	if value == nil {
		value = new(string)
		*value = defaultValue
	}
}

func MakeInQuery(num int) string {
	return " IN (?" + strings.Repeat(", ?", num-1) + ") "
}

func ParseSlug(str string) *string {
	regexpNonAuthorizedChars := regexp.MustCompile("[^a-zA-Z0-9-_]")
	regexpMultipleDashes := regexp.MustCompile("-+")

	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	slug := regexpNonAuthorizedChars.ReplaceAllString(str, "-")
	slug = regexpMultipleDashes.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-_")
	return &slug
}

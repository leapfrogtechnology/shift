package utils

import (
	"fmt"
	"net/url"
	"regexp"
)

// IsAlphaNumeric : Checks if the given value is alpha numeric
func IsAlphaNumeric(val string) bool {
	if CheckIsEmptyOrNull(val) {
		return false
	}
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	res := re.MatchString(val)
	fmt.Println(res)
	return res
}

// IsValidSlackWebHookURL : Check if the given value is url as well as slack web hook url
func IsValidSlackWebHookURL(urlVal string) bool {
	if CheckIsEmptyOrNull(urlVal) {
		return false
	}
	_, err := url.ParseRequestURI(urlVal)
	if err != nil {
		return false
	}
	return true
}

// CheckIsEmptyOrNull : Check if the given value is empty or not
func CheckIsEmptyOrNull(val string) bool {
	if len(val) == 0 {
		return true
	}
	return false
}

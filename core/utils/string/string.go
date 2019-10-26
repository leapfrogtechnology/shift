package utils

import (
	"fmt"
	"net/url"
	"regexp"
)

// IsAlphaNumeric : Checks if the given value is alpha numeric
func IsAlphaNumeric(val string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	res := re.MatchString(val)
	fmt.Println(res)
	return res
}

// IsValidSlackWebHookURL : Check if the given value is url as well as slack web hook url
func IsValidSlackWebHookURL(urlVal string) bool {
	_, err := url.ParseRequestURI(urlVal)
	if err != nil {
		return false
	}
	return true
}

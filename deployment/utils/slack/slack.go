package slack

import (
	"fmt"
	"strings"

	"github.com/leapfrogtechnology/shift/deployment/utils/http"
)

// Notify sends data to slack webhook url
func Notify(url string, text string, color string) {

	message := strings.Replace(text, "\"", "'", -1)

	var content string

	if strings.Contains(message, "error") {

		content = fmt.Sprintf(`{"attachments": [{"text": %q, "color": "%s"}]}`, "```"+message+"```", color)
	} else {

		content = fmt.Sprintf(`{"attachments": [{"text": %q, "color": "%s"}]}`, message, color)
	}

	http.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(content)).
		Post(url)
}

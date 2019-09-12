package slack

import (
	"fmt"

	"github.com/leapfrogtechnology/shift/deployment/utils/http"
)

// Notify sends data to slack webhook url
func Notify(url string, text string, color string) {
	http.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(fmt.Sprintf(`{"attachments": [{"text": "%s", "color": "%s"}]}`, text, color))).
		Post(url)
}

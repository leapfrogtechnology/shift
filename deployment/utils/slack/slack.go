package slack

import (
	"fmt"
	"strings"

	"github.com/leapfrogtechnology/shift/deployment/utils/http"
)

// Notify sends data to slack webhook url
func Notify(url string, text string, color string) {

	message := strings.Replace(text, "\"", "'", -1)

	content := fmt.Sprintf(`{"attachments": [{"text": %q, "color": "%s"}]}`, "@here\n"+message, color)

	resp, _ := http.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(content)).
		Post(url)

	fmt.Println("Slack Hook:- " + resp.String())
}

package terraform

import (
	"os"

	"github.com/leapfrogtechnology/shift/core/utils/http"
)

// ActivateLocalRun sends requests to terraform to set local run as default.
func ActivateLocalRun(workspace string) {
	http.Client.R().
		SetHeader("Authorization", "Bearer "+os.Getenv("TERRAFORM_TOKEN")).
		SetHeader("Content-Type", "application/vnd.api+json").
		SetBody([]byte("{\"data\": {\"type\": \"workspaces\", \"attributes\": {\"operations\": false}}}")).
		Patch("https://app.terraform.io/api/v2/organizations/lftechnology/workspaces/" + workspace)
}

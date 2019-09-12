package deploy

import (
	"encoding/json"
	"os/user"

	"github.com/leapfrogtechnology/shift/cli/services/mq/trigger"
)

// TriggerRequest defined the structure for Trigger
type TriggerRequest struct {
	Project    string `json:"project"`
	Deployment string `json:"deployment"`
	User       string `json:"user"`
}

// Run triggers the deployment.
func Run(project string, deployment string) {
	user, _ := user.Current()

	triggerRequest := TriggerRequest{
		Project:    project,
		Deployment: deployment,
		User:       user.Name,
	}

	triggerRequestJSON, _ := json.Marshal(triggerRequest)

	trigger.Publish(triggerRequestJSON)
}

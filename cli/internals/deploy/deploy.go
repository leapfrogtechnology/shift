package deploy

import (
	"encoding/json"

	"github.com/leapfrogtechnology/shift/cli/services/mq/trigger"
)

// TriggerRequest defined the structure for Trigger
type TriggerRequest struct {
	Project    string `json:"project"`
	Deployment string `json:"deployment"`
}

// Run triggers the deployment.
func Run(project string, deployment string) {
	triggerRequest := TriggerRequest{
		Project:    project,
		Deployment: deployment,
	}

	triggerRequestJSON, _ := json.Marshal(triggerRequest)

	trigger.Publish(triggerRequestJSON)
}

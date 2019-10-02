package project

type deployment struct {
	Name            string `json:"name"`
	Platform        string `json:"platform"`
	Type            string `json:"type"`
	BuildCommand    string `json:"buildCommand"`
	DistFolder      string `json:"distFolder"`
	Port            string `json:"port"`
	HealthCheckPath string `json:"healthCheckPath"`
	DockerFilePath  string `json:"dockerFilePath"`
	SlackURL        string `json:"slackURL"`
}

type terraformOutput struct {
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type infrastructure struct {
	BucketName                 terraformOutput `json:"bucketName"`
	FrontendWebURL             terraformOutput `json:"appUrl"`
	BackendClusterName         terraformOutput `json:"backendClusterName"`
	BackendContainerDefinition terraformOutput `json:"backendContainerDefinition"`
	BackendServiceID           terraformOutput `json:"backendServiceId"`
	BackendTaskDefinitionID    terraformOutput `json:"backendTaskDefinitionId"`
	RepoURL                    terraformOutput `json:"repoUrl"`
}

// TriggerRequest defines the requst to trigger a deployment.
type TriggerRequest struct {
	Project    string `json:"project"`
	Deployment string `json:"deployment"`
	User       string `json:"user"`
}

// Project defines the fields required for a single project.
type Project struct {
	ProjectName string         `json:"projectName"`
	Deployment  deployment     `json:"deployment"`
	Data        infrastructure `json:"data"`
}

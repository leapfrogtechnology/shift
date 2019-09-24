package project

type deployment struct {
	Name            string `json:"name"`
	Platform        string `json:"platform"`
	AccessKey       string `json:"accessKey"`
	SecretKey       string `json:"secretKey"`
	Type            string `json:"type"`
	GitProvider     string `json:"gitProvider"`
	GitToken        string `json:"gitToken"`
	CloneURL        string `json:"cloneUrl"`
	BuildCommand    string `json:"buildCommand"`
	DistFolder      string `json:"distFolder"`
	Port            string `json:"port"`
	HealthCheckPath string `json:"healthCheckPath"`
	DockerFilePath  string `json:"dockerFilePath"`
	SlackURL        string `json:"slackURL"`
	RepoName        string `json:"repoName"`
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
	BackendServiceId           terraformOutput `json:"backendServiceId"`
	BackendTaskDefinitionId    terraformOutput `json:"backendTaskDefinitionId"`
	RepoUrl                    terraformOutput `json:"repoUrl"`
}

// Response defines the response from message queue.
type Response struct {
	ProjectName string         `json:"projectName"`
	Deployment  deployment     `json:"deployment"`
	Data        infrastructure `json:"data"`
}

// TriggerRequest defines the requst to trigger a deployment.
type TriggerRequest struct {
	Project    string `json:"project"`
	Deployment string `json:"deployment"`
	User       string `json:"user"`
}

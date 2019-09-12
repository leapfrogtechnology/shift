package utils

type deployment struct {
	Name            string `json:"name"`
	Platform        string `json:"platform"`
	AccessKey       string `json:"accessKey"`
	SecretKey       string `json:"secretKey"`
	Type            string `json:"type"`
	GitProvider     string `json:"gitProvider"`
	GitToken        string `json:"gitToken"`
	CloneUrl        string `json:"cloneUrl"`
	BuildCommand    string `json:"buildCommand"`
	DistFolder      string `json:"distFolder"`
	Port            string `json:"port"`
	HealthCheckPath string `json:"healthCheckPath"`
	DockerFilePath  string `json:"dockerFilePath"`
	SlackURL        string `json:"slackURL"`
}
type Client struct {
	Project    string     `json:"projectName"`
	Deployment deployment `json:"deployment"`
}

type FrontendResult struct {
	Project    string                  `json:"projectName"`
	Deployment deployment              `json:"deployment"`
	Data       FrontendTerraformOutput `json:"data"`
}

type terraformOutput struct {
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type BackendTerraformOutput struct {
	BackendClusterName         terraformOutput `json:"backendClusterName"`
	BackendContainerDefinition terraformOutput `json:"backendContainerDefinition"`
	BackendServiceId           terraformOutput `json:"backendServiceId"`
	BackendTaskDefinitionId    terraformOutput `json:"backendTaskDefinitionId"`
	BackendUrl                 terraformOutput `json:"appUrl"`
	RepoUrl                    terraformOutput `json:"repoUrl"`
}

type BackendResult struct {
	Project    string                 `json:"projectName"`
	Deployment deployment             `json:"deployment"`
	Data       BackendTerraformOutput `json:"data"`
}

type FrontendTerraformOutput struct {
	BucketName     terraformOutput `json:"bucketName"`
	FrontendWebUrl terraformOutput `json:"appUrl"`
}

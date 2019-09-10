package project

type deployment struct {
	Name         string `json:"name"`
	Platform     string `json:"platform"`
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	Type         string `json:"type"`
	GitProvider  string `json:"gitProvider"`
	GitToken     string `json:"gitToken"`
	CloneURL     string `json:"cloneURL"`
	BuildCommand string `json:"buildCommand"`
	DistFolder   string `json:"distFolder"`
	Bucket       string `json:"bucket"`
}

type infrastructure struct {
	BucketName string `json:"bucketName"`
	URL        string `json:"url"`
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
}

package utils

type deployment struct {
	Name         string `json:"name"`
	Platform     string `json:"platform"`
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	Type         string `json:"type"`
	GitProvider  string `json:"gitProvider"`
	GitToken     string `json:"gitToken"`
	CloneUrl     string `json:"cloneUrl"`
	BuildCommand string `json:"buildCommand"`
	DistFolder   string `json:"distFolder"`
}
type Client struct {
	Project    string     `json:"projectName"`
	Deployment deployment `json:"deployment"`
}
type Frontend struct {
	BucketName string `json:"bucketName"`
	Url        string `json:"url"`
}
type FrontendResult struct {
	Project    string     `json:"projectName"`
	Deployment deployment `json:"deployment"`
	Data   Frontend `json:"data"`
}

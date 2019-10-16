package structs

// Env defines the structure for a single environment
type Env struct {
	Bucket  string `json:"bucket"`
	Cluster string `json:"cluster"`
}

// Project defines the overall structure for a project deployment.
type Project struct {
	Name            string         `json:"name"`
	Platform        string         `json:"platform"`
	Profile         string         `json:"profile"`
	Region          string         `json:"region"`
	Type            string         `json:"type"`
	DistDir         string         `json:"distDir"`
	SlackURL        string         `json:"slackURL"`
	Port            string         `json:"port"`
	DockerFilePath  string         `json:"dockerFilePath"`
	HealthCheckPath string         `json:"healthCheckPath"`
	Env             map[string]Env `json:"env"`
}

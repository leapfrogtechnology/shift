package structs

// Env defines the structure for a single environment
type Env struct {
	Bucket       string `json:"bucket,omitempty"`
	Cluster      string `json:"cluster,omitempty"`
	BuildCommand string `json:"buildCommand,omitempty"`
}

// Project defines the overall structure for a project deployment.
type Project struct {
	Name            string         `json:"name,omitempty"`
	Platform        string         `json:"platform,omitempty"`
	Profile         string         `json:"profile,omitempty"`
	Region          string         `json:"region,omitempty"`
	Type            string         `json:"type,omitempty"`
	DistDir         string         `json:"distDir,omitempty"`
	SlackURL        string         `json:"slackURL,omitempty"`
	Port            string         `json:"port,omitempty"`
	DockerFilePath  string         `json:"dockerFilePath,omitempty"`
	HealthCheckPath string         `json:"healthCheckPath,omitempty"`
	Env             map[string]Env `json:"env"`
}

// Frontend defined the output for frontend infrastructure
type Frontend struct {
	Bucket string `json:"bucket"`
	URL    string `json:"url"`
}

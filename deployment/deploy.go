package main

import (
	"encoding/json"

	"github.com/leapfrogtechnology/shift/deployment/internals/frontend"
	"github.com/leapfrogtechnology/shift/deployment/services/aws/s3"
	"github.com/leapfrogtechnology/shift/deployment/services/mq"
)

type deployment struct {
	Name         string `json:"name"`
	Platform     string `json:"platform"`
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	Type         string `json:"type"`
	GitProvider  string `json:"gitProvider"`
	GitToken     string `json:"gitToken"`
	CloneURL     string `json:"cloneURL"`
	BuildCommand string `json:"BuildCommand"`
	DistFolder   string `json:"distFolder"`
	Bucket       string `json:"bucket"`
}

type infrastructure struct {
	BucketName string `json:"bucketName"`
	URL        string `json:"url"`
}

type projectResponse struct {
	ProjectName string     `json:"projectName"`
	Deployment  deployment `json:"deployment"`
	Data        infrastructure
}

func deploy(msg []byte) {
	project := projectResponse{}
	json.Unmarshal(msg, &project)

	buildData := frontend.BuildData{
		GitToken:     project.Deployment.GitToken,
		Platform:     project.Deployment.Platform,
		CloneURL:     project.Deployment.CloneURL,
		BuildCommand: project.Deployment.BuildCommand,
		DistFolder:   project.Deployment.DistFolder,
		AccessKey:    project.Deployment.AccessKey,
		SecretKey:    project.Deployment.SecretKey,
	}

	frontend.Build(buildData)

	s3.Deploy(s3.Data{
		AccessKey:  project.Deployment.AccessKey,
		SecretKey:  project.Deployment.SecretKey,
		Bucket:     project.Data.BucketName,
		URL:        project.Data.URL,
		DistFolder: project.Deployment.DistFolder,
	})
}

func main() {
	mq.Consume(deploy)
}

package frontend

import (
	"fmt"

	"github.com/leapfrogtechnology/shift/deployment/utils/shell"
	"github.com/leapfrogtechnology/shift/deployment/utils/spinner"
)

// BuildData provides data for build.
type BuildData struct {
	GitToken     string `json:"gitToken"`
	Platform     string `json:"platform"`
	CloneURL     string `json:"cloneURL"`
	BuildCommand string `json:"buildCommand"`
	DistFolder   string `json:"distFolder"`
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
}

// Build executes build command
func Build(data BuildData) error {
	spinner.Start("")
	defer spinner.Stop()

	fmt.Println(data.CloneURL)
	cloneURL := "https://" + data.GitToken + "@" + data.CloneURL[8:]

	err := shell.Execute("export AWS_ACCESS_KEY_ID=" + data.AccessKey + " && export AWS_SECRET_ACCESS_KEY=" + data.SecretKey + " && cd /tmp && rm -rf artifact && mkdir artifact && cd artifact &&" + "git clone " + cloneURL + " source && cd source && " + data.BuildCommand)

	return err
}

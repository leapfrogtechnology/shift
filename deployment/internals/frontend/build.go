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
}

// Build executes build command
func Build(data BuildData) {
	spinner.Start("")
	defer spinner.Stop()

	fmt.Println(data.CloneURL)
	cloneURL := "https://" + data.GitToken + "@" + data.CloneURL[8:]

	shell.Execute("rm -rf artifact && mkdir artifact && cd artifact &&" + "git clone " + cloneURL + "&&" + data.BuildCommand)
}

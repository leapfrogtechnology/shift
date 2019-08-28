package utils

import (
	"bytes"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora"
	"os/exec"
	"time"
)

func initTerraform(workspaceDir string) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)  // Build our new spinner
	s.Prefix = "  "
	s.Suffix = "  Initializing"
	_ = s.Color("cyan", "bold")
	s.Start()
	cmd := exec.Command("terraform", "init")
	cmd.Dir = workspaceDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ":" + stderr.String())
		return
	}
	s.Stop()
}

func planTerraform(workspaceDir string) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)  // Build our new spinner
	s.Prefix = "  "
	s.Suffix = "  Planning"
	_ = s.Color("cyan", "bold")
	s.Start()
	cmd := exec.Command("terraform", "plan")
	cmd.Dir = workspaceDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ":" + stderr.String())
		return
	}
	s.Stop()
}

func applyTerraform(workspaceDir string) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)  // Build our new spinner
	s.Prefix = "  "
	s.Suffix = "  Applying Changes"
	_ = s.Color("cyan", "bold")
	s.Start()
	cmd := exec.Command("terraform", "apply", "--auto-approve")
	cmd.Dir = workspaceDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ":" + stderr.String())
		return
	}
	s.Stop()
}

func getTerraformOutput(workspaceDir string) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)  // Build our new spinner
	s.Prefix = "  "
	s.Suffix = "  Generating Output"
	_ = s.Color("cyan", "bold")
	s.Start()
	cmd := exec.Command("terraform", "output", "-json")
	cmd.Dir = workspaceDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ":" + stderr.String())
		return
	}
	s.Stop()
	fmt.Println(aurora.Cyan(stdout.String()))
}

func RunInfrastructureChanges(workspaceDir string)  {
	initTerraform(workspaceDir)
	planTerraform(workspaceDir)
	applyTerraform(workspaceDir)
	getTerraformOutput(workspaceDir)
}

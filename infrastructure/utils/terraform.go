package utils

import (
	"bytes"
	"fmt"
	"github.com/briandowns/spinner"
	"os/exec"
	"time"
)

func initTerraform(workspaceDir string) error {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
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
		fmt.Println()
		s.Stop()
		return err
	}
	s.Stop()
	return nil
}

func planTerraform(workspaceDir string) error {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
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
		LogError(err, stderr.String())
		return err
	}
	s.Stop()
	return nil
}

func applyTerraform(workspaceDir string) error {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
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
		LogError(err, stderr.String())
		return err
	}
	s.Stop()
	return nil
}

func getTerraformOutput(workspaceDir string) (string, error) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
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
		LogError(err, stderr.String())
		s.Stop()
		return "", err
	}
	s.Stop()
	return stdout.String(), nil
}

func RunInfrastructureChanges(workspaceDir string) (string, error) {
	LogInfo("Initializing")
	err := initTerraform(workspaceDir)
	if err != nil {
		LogError(err, "Something Went Wrong")
		return "", err
	}
	LogInfo("Planning")
	err = planTerraform(workspaceDir)
	if err != nil {
		LogError(err, "Something Went Wrong")
		return "", err
	}
	LogInfo("Applying")
	err = applyTerraform(workspaceDir)
	if err != nil {
		LogError(err, "Something Went Wrong")
		return "", err
	}
	infrastructureInfo, err := getTerraformOutput(workspaceDir)
	if err != nil {
		LogError(err, "Something Went Wrong")
		return "", err
	}
	return infrastructureInfo, err
}

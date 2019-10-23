package terraform

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/infrastructure/services/terraform"
)

// TODO use terraform Library to remove the dependency of installing terraform
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
		logger.LogError(err, stderr.String())
		s.Stop()
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
		logger.LogError(err, stderr.String())
		s.Stop()
		return "", err
	}
	s.Stop()

	return stdout.String(), err
}

// RunInfrastructureChanges starts terraform.
func RunInfrastructureChanges(workspaceDir string, workspaceName string) (string, error) {
	logger.LogInfo("Initializing")
	err := initTerraform(workspaceDir)

	// Set local execution instead of remote.
	terraform.ActivateLocalRun(workspaceName)

	if err != nil {
		logger.LogError(err, "Couldnot initialize")

		return "", err
	}

	logger.LogInfo("Applying")
	err = applyTerraform(workspaceDir)

	if err != nil {
		logger.LogError(err, "Failed to apply changes")

		return "", err
	}

	out, err := getTerraformOutput(workspaceDir)

	if err != nil {
		logger.LogError(err, "Failed to get terraform output")

		return "", err
	}

	return out, err
}

package terraform

import (
	"bytes"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"os/exec"
	"strings"
	"time"
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

func getTerraformOutputFrontend(workspaceDir string) (string, string, error) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
	s.Prefix = "  "
	s.Suffix = "  Generating Output"
	_ = s.Color("cyan", "bold")
	s.Start()
	cmd1 := exec.Command("terraform", "output", "bucket_name")
	cmd1.Dir = workspaceDir
	var stdout1 bytes.Buffer
	var stderr1 bytes.Buffer
	cmd1.Stdout = &stdout1
	cmd1.Stderr = &stderr1
	err := cmd1.Run()
	if err != nil {
		logger.LogError(err, stderr1.String())
		s.Stop()
		return "", "", err
	}
	bucketName := strings.TrimSuffix(stdout1.String(), "\n")
	cmd2 := exec.Command("terraform", "output", "frontend_web_url")
	cmd2.Dir = workspaceDir
	var stdout2 bytes.Buffer
	var stderr2 bytes.Buffer
	cmd2.Stdout = &stdout2
	cmd2.Stderr = &stderr2
	err = cmd2.Run()
	if err != nil {
		logger.LogError(err, stderr2.String())
		s.Stop()
		return "", "", err
	}
	url := "https://" + strings.TrimSuffix(stdout2.String(), "\n")
	s.Stop()

	return bucketName, url, err
}

func RunFrontendInfrastructureChanges(workspaceDir string) (string, string, error) {
	logger.LogInfo("Initializing")
	err := initTerraform(workspaceDir)
	if err != nil {
		logger.LogError(err, "Something Went Wrong")
		return "", "", err
	}
	logger.LogInfo("Applying")
	err = applyTerraform(workspaceDir)
	if err != nil {
		logger.LogError(err, "Something Went Wrong")
		return "", "", err
	}
	bucketName, url, err := getTerraformOutputFrontend(workspaceDir)
	if err != nil {
		logger.LogError(err, "Something Went Wrong")
		return "", "", err
	}
	return bucketName, url, err
}

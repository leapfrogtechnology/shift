package utils

import (
	"bytes"
	"fmt"
	"github.com/briandowns/spinner"
	"os/exec"
	"strings"
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
		s.Stop()
		return err
	}
	s.Stop()
	return nil
}

func getTerraformOutput(workspaceDir string) (string, string, error) {
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
		LogError(err, stderr1.String())
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
		LogError(err, stderr2.String())
		s.Stop()
		return "", "", err
	}
	url := "https://" + strings.TrimSuffix(stdout2.String(), "\n")
	s.Stop()

	return bucketName, url, err
}

func RunInfrastructureChanges(workspaceDir string) (string, string, error) {
	LogInfo("Initializing")
	err := initTerraform(workspaceDir)
	if err != nil {
		LogError(err, "Something Went Wrong")
		return "", "", err
	}
	//LogInfo("Planning")
	//err = planTerraform(workspaceDir)
	//if err != nil {
	//	LogError(err, "Something Went Wrong")
	//	return "", "", err
	//}
	LogInfo("Applying")
	err = applyTerraform(workspaceDir)
	if err != nil {
		LogError(err, "Something Went Wrong")
		return "", "", err
	}
	bucketName, url, err := getTerraformOutput(workspaceDir)
	if err != nil {
		LogError(err, "Something Went Wrong")
		return "", "", err
	}
	return bucketName, url, err
}

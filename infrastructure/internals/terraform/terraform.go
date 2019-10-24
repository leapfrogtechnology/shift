package terraform

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/leapfrogtechnology/shift/core/structs"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/core/utils/spinner"
	"github.com/leapfrogtechnology/shift/infrastructure/internals/terraform/templates/providers/aws/template"
	"github.com/leapfrogtechnology/shift/infrastructure/services/terraform"
)

// TODO use terraform Library to remove the dependency of installing terraform

func initTerraform(workspaceDir string) error {
	logger.LogInfo("  Initializing")
	spinner.Start(" ")

	cmd := exec.Command("terraform", "init")
	cmd.Dir = workspaceDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		spinner.Stop()
		return err
	}
	spinner.Stop()
	return nil
}

func applyTerraform(workspaceDir string) error {

	logger.LogInfo("  Applying Changes")
	spinner.Start(" ")

	cmd := exec.Command("terraform", "apply", "--auto-approve")
	cmd.Dir = workspaceDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		logger.LogError(err, "Error While terraform apply")
		spinner.Stop()
		return err
	}
	spinner.Stop()
	return nil
}

func getTerraformOutput(workspaceDir string) (string, error) {

	logger.LogInfo("  Generating Output")
	spinner.Start(" ")
	cmd := exec.Command("terraform", "output", "-json")
	cmd.Dir = workspaceDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logger.LogError(err, stderr.String())
		spinner.Stop()
		return "", err
	}
	spinner.Stop()

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

// DestroyInfrastructure destroys existing infrastructure
func DestroyInfrastructure(workspaceDir string) error {
	logger.LogInfo("Distroying Infrastructure...")
	spinner.Start(" ")
	cmd := exec.Command("terraform", "destroy", "--auto-approve")
	cmd.Dir = workspaceDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	message := " ERROR While Destroying Infrastructure"
	if err != nil {
		logger.LogError(err, message)
		spinner.Stop()
		return errors.New(err.Error() + message)
	}

	spinner.Stop()
	return nil

}

// MakeTempAndDestroy create infrastructure template and distroy the infrastructure
func MakeTempAndDestroy(project structs.Project, environment, workspaceDir string) error {

	logger.LogInfo("Generating Templates....")
	if project.Type == "Frontend" {
		template.GenerateFrontendTemplate(project, workspaceDir, environment)
	} else {
		// TODO backend
	}
	initTerraform(workspaceDir)
	err := DestroyInfrastructure(workspaceDir)
	return err
}

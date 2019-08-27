package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type s3WebsiteVariables struct {
	CLIENT_WORKSPACE   string `json:"client_workspace"`
	AWS_REGION         string `json:"aws_region"`
	AWS_ACCESS_KEY     string `json:"aws_access_key"`
	AWS_SECRET_KEY     string `json:"aws_secret_key"`
	AWS_S3_BUCKET_NAME string `json:"aws_s3_bucket_name"`
	TERRAFORM_TOKEN    string `json:"terraform_token"`
}

func generateTemplateFile(templateLocation string, s3Args s3WebsiteVariables, terraformPath string) {
	tpl, err := pongo2.FromFile(templateLocation)
	if err != nil {
		panic(err)
	}
	out, err := tpl.Execute(pongo2.Context{"info": s3Args})
	if err != nil {
		panic(err)
	}
	terraformFileName := terraformPath + "/infrastructure.tf"
	err = os.MkdirAll(terraformPath, 0700)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(terraformFileName, []byte(out), 0600)
	if err != nil {
		panic(err)
	}
}

func initTerraform(workspaceDir string) {
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
	fmt.Println("Initialized")
}

func planTerraform(workspaceDir string) {
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
	fmt.Println("Planned")
}

func applyTerraform(workspaceDir string) {
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
	fmt.Println("Applied Changes")
}

func getTerraformOutput(workspaceDir string) {
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
	fmt.Println(stdout.String())
}

func main() {
	credentialsJsonFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer credentialsJsonFile.Close()
	byteValue, _ := ioutil.ReadAll(credentialsJsonFile)
	var s3Args s3WebsiteVariables
	err = json.Unmarshal(byteValue, &s3Args)
	if err != nil {
		panic(err)
	}
	workspaceDir := filepath.Join("/tmp", s3Args.CLIENT_WORKSPACE)
	generateTemplateFile("templates/providers/aws/s3-website/terraform.tpl", s3Args, workspaceDir)
	initTerraform(workspaceDir)
	planTerraform(workspaceDir)
	applyTerraform(workspaceDir)
	getTerraformOutput(workspaceDir)
}

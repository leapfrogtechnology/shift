package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

const accessKey = ""
const secretKey = ""
const repoUrl = ""

func main() {
	_ = os.Setenv("AWS_ACCESS_KEY_ID", accessKey)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", secretKey)
	_ = os.Setenv("REGION_NAME", "us-east-1")
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	command := fmt.Sprintf("$(AWS_ACCESS_KEY_ID=%s AWS_SECRET_ACCESS_KEY=%s aws ecr get-login --no-include-email --region %s) && docker build docker -t %s && docker push %s", accessKey, secretKey, "us-east-1", repoUrl, repoUrl)
	fmt.Println(command)
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Dir = path
	fmt.Println(cmd.Env)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		panic(err)
	} else {
		fmt.Println(stdout.String())
	}
}

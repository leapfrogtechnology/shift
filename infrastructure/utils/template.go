package utils

import (
	"github.com/flosch/pongo2"
	"io/ioutil"
	"os"
)

type FrontendInfrastructureVariables struct {
	CLIENT_NAME        string `json:"client_name"`
	AWS_REGION         string `json:"aws_region"`
	AWS_ACCESS_KEY     string `json:"aws_access_key"`
	AWS_SECRET_KEY     string `json:"aws_secret_key"`
	AWS_S3_BUCKET_NAME string `json:"aws_s3_bucket_name"`
	TERRAFORM_TOKEN    string `json:"terraform_token"`
}

func GenerateFrontendTemplateFile(template string, s3Args FrontendInfrastructureVariables, terraformPath string) {
	tpl, err := pongo2.FromString(template)
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

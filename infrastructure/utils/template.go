package utils

import (
	"github.com/flosch/pongo2"
	"io/ioutil"
	"os"
)

type infrastructure struct {
	Client Client
	Token  string
}

func GenerateFrontendTemplateFile(template string, client Client, terraformPath string) error {
	token := os.Getenv("TERRAFORM_TOKEN")
	infrastructure := infrastructure{client, token}

	tpl, err := pongo2.FromString(template)
	if err != nil {
		return err
	}
	out, err := tpl.Execute(pongo2.Context{"info": infrastructure})
	if err != nil {
		return err
	}
	terraformFileName := terraformPath + "/infrastructure.tf"
	err = os.MkdirAll(terraformPath, 0700)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(terraformFileName, []byte(out), 0600)
	if err != nil {
		return err
	}
	return nil
}

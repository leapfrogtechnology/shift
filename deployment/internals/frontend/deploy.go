package frontend

import (
	"github.com/leapfrogtechnology/shift/core/structs"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/core/utils/shell"
	"github.com/leapfrogtechnology/shift/deployment/services/platforms/aws/s3"
)

// Deploy uploads files to S3 bucket.
func Deploy(project structs.Project, environment string) {
	buildCommand := project.Env[environment].BuildCommand

	if buildCommand != "" {
		shell.Execute(buildCommand)
	}

	err := s3.Deploy(s3.Data{
		Profile: project.Profile,
		Region:  project.Region,
		Bucket:  project.Env[environment].Bucket,
		DistDir: project.DistDir,
	})

	logger.FailOnError(err, "Deployment Failed.")
}

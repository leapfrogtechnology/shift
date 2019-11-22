package backend

import (
	"github.com/leapfrogtechnology/shift/core/structs"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/core/utils/shell"
)

func buildImage(dockerFilePath string, image string, tag string) {
	err := shell.Execute("docker build " + dockerFilePath + " -t " + image + ":" + tag)
	logger.FailOnError(err, "Failed to build docker image.")
}

func loginECR(region string, profile string) {
	err := shell.Execute("$(aws ecr get-login --no-include-email --region " + region + " --profile " + profile + ")")
	logger.FailOnError(err, "Failed to login to ECR.")
}

func pushImage(image string, tag string) {
	err := shell.Execute("docker push " + image + ":" + tag)
	logger.FailOnError(err, "Failed to push image to ECR.")
}

func updateService(region string, profile string, cluster string, service string) {
	err := shell.Execute("aws ecs update-service --cluster " + cluster + " --service " + service + " --force-new-deployment --region " + region + " --profile " + profile)
	logger.FailOnError(err, "Failed to update service.")
}

// Deploy uploads a new Docker image and updates the service.
func Deploy(project structs.Project, environment string) {
	env := project.Env[environment]

	image := env.Image
	cluster := env.Cluster
	service := env.Service

	logger.Info("Building Image:")
	buildImage(project.DockerFilePath, image, environment)

	logger.Info("Logging in to ECR:")
	loginECR(project.Region, project.Profile)

	logger.Info("Pushing Docker Image to ECR:")
	pushImage(image, environment)

	logger.Info("Updating Service:")
	updateService(project.Region, project.Profile, cluster, service)

	logger.Info("Deployment Successfull. 🎉 🎉 🎉")
}

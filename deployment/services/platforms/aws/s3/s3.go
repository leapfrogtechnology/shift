package s3

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	awsService "github.com/leapfrogtechnology/shift/core/services/platforms/aws"
	fileUtil "github.com/leapfrogtechnology/shift/core/utils/file"
	"github.com/leapfrogtechnology/shift/core/utils/logger"
	"github.com/leapfrogtechnology/shift/core/utils/spinner"
)

// Data contains the data needed to deploy to S3 bucket
type Data struct {
	Profile string
	Region  string
	Bucket  string
	DistDir string
}

// Deploy to S3 bucket
func Deploy(data Data) error {
	session := awsService.GetSession(data.Profile, data.Region)
	uploader := s3manager.NewUploader(session)

	fileList := []string{}

	filepath.Walk(data.DistDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			fileList = append(fileList, path)

			return nil
		})

	spinner.Start("Uploading")

	for _, file := range fileList {
		f, _ := os.Open(file)

		logger.Log("Uploading: " + file)

		key := strings.TrimPrefix(file, data.DistDir)
		contentType := fileUtil.GetFileContentType(file)

		// Upload the file to S3.
		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(data.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
			Body:        f,
		})

		if err != nil {
			return err
		}

		logger.Success("Success")
	}

	logger.Info("🎉 🎉 🎉  Files Uploaded Successfully. 🎉 🎉 🎉")

	spinner.Stop()

	return nil
}

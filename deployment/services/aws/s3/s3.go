package s3

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	awsService "github.com/leapfrogtechnology/shift/deployment/services/aws"
	fileUtil "github.com/leapfrogtechnology/shift/deployment/utils/file"
	"github.com/leapfrogtechnology/shift/deployment/utils/spinner"
)

// Data contains the data needed to deploy to S3 bucket
type Data struct {
	AccessKey  string
	SecretKey  string
	Bucket     string
	DistFolder string
	URL        string
}

// Deploy to S3 bucket
func Deploy(data Data) {
	session := awsService.GetSession(data.AccessKey, data.SecretKey)
	uploader := s3manager.NewUploader(session)

	fileList := []string{}

	filepath.Walk("./artifact/"+data.DistFolder,
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

		key := strings.TrimPrefix(file, "artifact/"+data.DistFolder)
		contentType := fileUtil.GetFileContentType(file)

		// Upload the file to S3.
		output, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(data.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
			Body:        f,
		})

		if err != nil {
			fmt.Println("error")
			fmt.Println(err)
		}

		fmt.Println(output.Location)
	}

	fmt.Println("Uploaded all files.")
	fmt.Println("Project deployed at " + data.URL)
	spinner.Stop()
}

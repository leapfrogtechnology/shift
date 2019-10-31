package s3

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	awsService "github.com/leapfrogtechnology/shift/core/services/platforms/aws"
	fileUtil "github.com/leapfrogtechnology/shift/core/utils/file"

	"github.com/schollz/progressbar/v2"
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

	bar := progressbar.NewOptions(len(fileList),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]▌▌[reset]",
			SaucerPadding: "[green]░[reset]",
			BarStart:      "╢",
			BarEnd:        "╟",
		}))

	for _, file := range fileList {
		f, _ := os.Open(file)

		key := strings.TrimPrefix(file, data.DistDir)
		contentType := fileUtil.GetFileContentType(file)

		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(data.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
			Body:        f,
		})

		if err != nil {
			return err
		}
		bar.Describe("Uploading... " + file)
		bar.Add(1)
	}

	fmt.Println("\n\n" + strconv.Itoa(len(fileList)) + " Files Uploaded Successfully. 🎉 🎉 🎉")

	return nil
}

// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package s3util

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type DownloadInput struct {
	Downloader *s3manager.Downloader
	Bucket     string
	Key        string
	Path       string
}

func Download(input *DownloadInput) error {
	outputFile, err := os.Create(input.Path)
	if err != nil {
		return fmt.Errorf("error creating destination file %q: %w", input.Path, err)
	}
	defer outputFile.Close()
	_, err = input.Downloader.Download(outputFile, &s3.GetObjectInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
	})
	if err != nil {
		return fmt.Errorf("error downloading file s3://%s/%s: %w", input.Bucket, input.Key, err)
	}
	return nil
}

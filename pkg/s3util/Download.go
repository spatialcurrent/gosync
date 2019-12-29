// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package s3util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type DownloadInput struct {
	Downloader *s3manager.Downloader
	Bucket     string
	Key        string
	Path       string
	Parents    bool
}

func Download(input *DownloadInput) error {
	if input.Downloader == nil {
		return errors.New("downloader is nil")
	}
	if input.Parents {
		err := os.MkdirAll(filepath.Dir(input.Path), 0755)
		if err != nil {
			return fmt.Errorf("error creating parent directories for %q: %w", input.Path, err)
		}
	}
	file, err := os.Create(input.Path)
	if err != nil {
		return fmt.Errorf("error creating destination file %q: %w", input.Path, err)
	}
	_, err = input.Downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
	})
	if err != nil {
		_ = file.Close() // silently close output file
		return fmt.Errorf("error downloading file s3://%s/%s: %w", input.Bucket, input.Key, err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("error closing file after downloading from AWS s3: %w", err)
	}
	return nil
}

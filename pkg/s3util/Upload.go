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
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"strings"
)

type UploadInput struct {
	Uploader *s3manager.Uploader
	Path     string
	Bucket   string
	Key      string
}

func Upload(input *UploadInput) error {
	if input.Uploader == nil {
		return fmt.Errorf("uploader is nil")
	}

	if strings.HasPrefix(input.Path, "~") {
		return fmt.Errorf("source path %q starts with \"~\"", input.Path)
	}

	if err := CheckBucket(input.Bucket); err != nil {
		return fmt.Errorf("destination bucket is invalid: %w", err)
	}

	if err := CheckKey(input.Key); err != nil {
		return fmt.Errorf("destination key is invalid: %w", err)
	}

	file, err := os.OpenFile(input.Path, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("error opening source file at %q: %w", input.Path, err)
	}
	defer file.Close()

	_, err = input.Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("error uploading file to AWS S3: %w", err)
	}

	return nil
}

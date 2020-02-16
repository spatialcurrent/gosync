// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package s3util

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/spatialcurrent/go-lazy/pkg/lazy"
)

type UploadInput struct {
	Context  context.Context
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

	reader := lazy.NewLazyReader(func() (io.Reader, error) {
		file, err := os.OpenFile(input.Path, os.O_RDONLY, 0)
		if err != nil {
			return nil, fmt.Errorf("error opening source file at %q: %w", input.Path, err)
		}
		return file, nil
	})

	_, err := input.Uploader.UploadWithContext(input.Context, &s3manager.UploadInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("error uploading file to AWS S3: %w", err)
	}

	err = reader.Close()
	if err != nil {
		return fmt.Errorf("error closing file after uploading to AWS s3: %w", err)
	}

	return nil
}

// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package s3util

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/spatialcurrent/go-lazy/pkg/lazy"
)

type DownloadInput struct {
	Context    context.Context
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
	writer := lazy.NewLazyWriterAt(func() (io.WriterAt, error) {
		file, err := os.Create(input.Path)
		if err != nil {
			return nil, fmt.Errorf("error creating destination file %q: %w", input.Path, err)
		}
		return file, nil
	})
	_, err := input.Downloader.DownloadWithContext(input.Context, writer, &s3.GetObjectInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
	})
	if err != nil {
		_ = writer.Close() // silently close output file
		return fmt.Errorf("error downloading file s3://%s/%s: %w", input.Bucket, input.Key, err)
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing file after downloading from AWS s3: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error closing file after downloading from AWS s3: %w", err)
	}
	return nil
}

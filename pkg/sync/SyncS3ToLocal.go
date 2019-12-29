// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/sync/errgroup"
	"path/filepath"
	"strings"

	"github.com/spatialcurrent/gosync/pkg/s3util"
)

type SyncS3ToLocalInput struct {
	Bucket      string
	KeyPrefix   string
	Destination string
	Client      *s3.S3
	Downloader  *s3manager.Downloader
	Parents     bool
	Limit       int
	Verbose     bool
}

func SyncS3ToLocal(input *SyncS3ToLocalInput) error {

	if err := s3util.CheckBucket(input.Bucket); err != nil {
		return fmt.Errorf("error with source bucket %q: %w", input.Bucket, err)
	}

	if err := s3util.CheckKeyPrefix(input.KeyPrefix); err != nil {
		return fmt.Errorf("error with source key prefix %q: %w", input.KeyPrefix, err)
	}

	if strings.HasPrefix(input.Destination, "~") {
		return fmt.Errorf("destination %q starts with \"~\"", input.Destination)
	}

	if input.Limit == 0 {
		return nil
	}

	it := s3util.NewIterator(&s3util.NewIteratorInput{
		Client: input.Client,
		Bucket: input.Bucket,
		Prefix: input.KeyPrefix,
	})

	i := 0
	var g errgroup.Group
	for {
		object, err := it.Next()
		if err != nil {
			break
		}
		key := aws.StringValue(object.Key)
		r, err := filepath.Rel(input.KeyPrefix, key)
		if err != nil {
			return fmt.Errorf(
				"error calculating relative path between key prefix (%q) and object key (%q): %w",
				input.KeyPrefix,
				key,
				err,
			)
		}
		destinationPath := filepath.Join(input.Destination, r)
		if input.Verbose {
			fmt.Println(fmt.Sprintf("[ %d ] : s3://%s/%s => file://%s", i+1, input.Bucket, key, destinationPath))
		}
		g.Go(func() error {
			err := s3util.Download(&s3util.DownloadInput{
				Downloader: input.Downloader,
				Bucket:     input.Bucket,
				Key:        aws.StringValue(object.Key),
				Path:       destinationPath,
				Parents:    input.Parents,
			})
			if err != nil {
				return fmt.Errorf("error downloading file: %w", err)
			}
			return nil
		})
		i++
		if input.Limit > 0 && i == input.Limit {
			break
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	if err := it.Error(); err != nil {
		return fmt.Errorf("error iterating over source s3://%s/%s: %w", input.Bucket, input.KeyPrefix, err)
	}

	return nil
}

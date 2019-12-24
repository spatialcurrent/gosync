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

func SyncS3ToLocal(bucket string, keyPrefix string, destination string, s3Client *s3.S3, downloader *s3manager.Downloader, verbose bool) error {

	if err := s3util.CheckBucket(bucket); err != nil {
		return fmt.Errorf("error with source bucket %q: %w:", bucket, err)
	}

	if err := s3util.CheckKeyPrefix(keyPrefix); err != nil {
		return fmt.Errorf("error with source key prefix %q: %w:", keyPrefix, err)
	}

	if strings.HasPrefix(destination, "~") {
		return fmt.Errorf("destination %q starts with \"~\"", destination)
	}

	it := s3util.NewIterator(&s3util.NewIteratorInput{
		Client: s3Client,
		Bucket: bucket,
		Prefix: keyPrefix,
	})

	var g errgroup.Group
	for {
		object, err := it.Next()
		if err != nil {
			break
		}
		key := aws.StringValue(object.Key)
		r, err := filepath.Rel(keyPrefix, key)
		if err != nil {
			return fmt.Errorf("error calculating relative path between key prefix (%q) and object key (%q): %w", keyPrefix, key, err)
		}
		destinationPath := filepath.Join(destination, r)
		if verbose {
			fmt.Println(fmt.Sprintf("s3://%s/%s => file://%s", bucket, key, destinationPath))
		}
		g.Go(func() error {
			err := s3util.Download(&s3util.DownloadInput{
				Downloader: downloader,
				Bucket:     bucket,
				Key:        aws.StringValue(object.Key),
				Path:       destinationPath,
			})
			if err != nil {
				return fmt.Errorf("error downloading file: %w", err)
			}
			return nil
		})

	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	if err := it.Error(); err != nil {
		return fmt.Errorf("error iterating over source s3://%s/%s: %w", bucket, keyPrefix, err)
	}

	return nil
}

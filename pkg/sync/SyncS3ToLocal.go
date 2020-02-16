// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/spatialcurrent/gosync/pkg/group"
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
	PoolSize    int
	StopOnError bool
	Verbose     bool
	Timeout     time.Duration
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

	g, err := group.New(input.PoolSize, input.Limit, input.StopOnError)
	if err != nil {
		return fmt.Errorf("error creating concurrent execution group: %w", err)
	}

	i := 0
	for {
		object, err := it.Next()
		if err != nil {
			break
		}

		key := aws.StringValue(object.Key)

		destinationPath := ""
		// if the path path the prefix does not contain any directory separators
		if !strings.Contains(key[len(input.KeyPrefix):], "/") {
			// Set destination path as a file within the destination directory
			destinationPath = filepath.Join(input.Destination, filepath.Base(key))
		} else {
			r, err := filepath.Rel(input.KeyPrefix, key)
			if err != nil {
				return fmt.Errorf(
					"error calculating relative path between key prefix (%q) and object key (%q): %w",
					input.KeyPrefix,
					key,
					err,
				)
			}
			destinationPath = filepath.Join(input.Destination, r)
		}
		index := i
		g.Go(func() error {
			if input.Verbose {
				fmt.Printf("[ %d ] : s3://%s/%s => file://%s\n", index+1, input.Bucket, key, destinationPath)
			}
			ctx := context.Background()
			if int(input.Timeout) > 0 {
				c, cancel := context.WithTimeout(ctx, input.Timeout)
				if err != nil {
					return fmt.Errorf("error creating timeout: %w", err)
				}
				ctx = c
				defer cancel()
			}
			err = s3util.Download(&s3util.DownloadInput{
				Context:    ctx,
				Downloader: input.Downloader,
				Bucket:     input.Bucket,
				Key:        aws.StringValue(object.Key),
				Path:       destinationPath,
				Parents:    input.Parents,
			})
			if err != nil {
				return fmt.Errorf("error downloading from \"%s/%s\" to %q: %w", input.Bucket, key, destinationPath, err)
			}
			return nil
		})
		i += 1
		if input.Limit > 0 && i == input.Limit {
			break
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	fmt.Println("Done downloading files.")

	if err := it.Error(); err != nil {
		return fmt.Errorf("error iterating over source s3://%s/%s: %w", input.Bucket, input.KeyPrefix, err)
	}

	return nil
}

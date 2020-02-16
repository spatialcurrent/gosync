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

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/spatialcurrent/gosync/pkg/group"
	"github.com/spatialcurrent/gosync/pkg/s3util"
)

type SyncLocalToS3Input struct {
	Source      string
	Bucket      string
	KeyPrefix   string
	Uploader    *s3manager.Uploader
	PoolSize    int
	StopOnError bool
	Limit       int
	Verbose     bool
	Timeout     time.Duration
}

func SyncLocalToS3(input *SyncLocalToS3Input) error {

	if strings.HasPrefix(input.Source, "~") {
		return fmt.Errorf("source %q starts with \"~\"", input.Source)
	}

	if strings.HasPrefix(input.Bucket, "~") {
		return fmt.Errorf("destination bucket %q starts with \"~\"", input.Bucket)
	}

	if strings.HasPrefix(input.KeyPrefix, "~") {
		return fmt.Errorf("destination key prefix %q starts with \"~\"", input.KeyPrefix)
	}

	sourcePaths, err := CollectFiles([]string{input.Source})
	if err != nil {
		return fmt.Errorf("error collecting files from %q: %w", input.Source, err)
	}

	if strings.Contains(input.Bucket, WildcardCharacter) {
		return fmt.Errorf("destination bucket cannot contain wildcard: %q", input.Bucket)
	}

	if strings.Contains(input.KeyPrefix, WildcardCharacter) {
		return fmt.Errorf("destination key prefix cannot contain wildcard: %q", input.KeyPrefix)
	}

	sourceMaxLength := maxLength(sourcePaths)

	g, err := group.New(input.PoolSize, input.Limit, input.StopOnError)
	if err != nil {
		return fmt.Errorf("error creating concurrent execution group: %w", err)
	}
	for i, p := range sourcePaths {
		i := i
		p := p
		r, err := filepath.Rel(input.Source, p)
		if err != nil {
			return fmt.Errorf("error calculating relative path between %q and %q: %w", input.Source, p, err)
		}
		key := filepath.Join(input.KeyPrefix, r)
		g.Go(func() error {
			if input.Verbose {
				fmt.Printf("[ %d ] : %s => s3://%s/%s\n", i+1, fillRight(p, sourceMaxLength), input.Bucket, key)
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
			err := s3util.Upload(&s3util.UploadInput{
				Context:  ctx,
				Uploader: input.Uploader,
				Path:     p,
				Bucket:   input.Bucket,
				Key:      key,
			})
			if err != nil {
				return fmt.Errorf("error uploading from %q to \"%s/%s\": %w", p, input.Bucket, key, err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	return nil
}

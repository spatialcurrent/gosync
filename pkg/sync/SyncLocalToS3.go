// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/sync/errgroup"
	"path/filepath"
	"strings"

	"github.com/spatialcurrent/gosync/pkg/s3util"
)

func SyncLocalToS3(source string, bucket string, keyPrefix string, uploader *s3manager.Uploader, verbose bool) error {

	if strings.HasPrefix(source, "~") {
		return fmt.Errorf("source %q starts with \"~\"", source)
	}

	if strings.HasPrefix(bucket, "~") {
		return fmt.Errorf("destination bucket %q starts with \"~\"", bucket)
	}

	if strings.HasPrefix(keyPrefix, "~") {
		return fmt.Errorf("destination key prefix %q starts with \"~\"", keyPrefix)
	}

	sourcePaths, err := CollectFiles([]string{source})
	if err != nil {
		return fmt.Errorf("error collecting files from %q: %w", source, err)
	}

	if strings.Contains(bucket, WildcardCharacter) {
		return fmt.Errorf("destination bucket cannot contain wildcard: %q", bucket)
	}

	if strings.Contains(keyPrefix, WildcardCharacter) {
		return fmt.Errorf("destination key prefix cannot contain wildcard: %q", keyPrefix)
	}

	sourceMaxLength := maxLength(sourcePaths)

	var g errgroup.Group
	for i, p := range sourcePaths {
		p := p
		r, err := filepath.Rel(source, p)
		if err != nil {
			return fmt.Errorf("error calculating relative path between %q and %q: %w", source, p, err)
		}
		key := filepath.Join(keyPrefix, r)
		if verbose {
			fmt.Println(fmt.Sprintf("[ %d ] : %s => s3://%s/%s", i+1, fillRight(p, sourceMaxLength), bucket, key))
		}
		g.Go(func() error {
			err := s3util.Upload(&s3util.UploadInput{
				Uploader: uploader,
				Path:     p,
				Bucket:   bucket,
				Key:      key,
			})
			if err != nil {
				return fmt.Errorf("error uploading %q to \"%s/%s\": %w", p, bucket, key, err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	return nil
}

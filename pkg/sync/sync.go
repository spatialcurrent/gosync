// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spatialcurrent/gosync/pkg/awsutil"
	"github.com/spatialcurrent/gosync/pkg/s3util"
	"strings"
)

const (
	WildcardCharacter = "*"
)

func maxLength(s []string) int {
	max := 0
	for _, v := range s {
		if i := len(v); i > max {
			max = i
		}
	}
	return max
}

func fillRight(s string, l int) string {
	if len(s) < l {
		return s + strings.Repeat(" ", l-len(s))
	}
	return s
}

func splitUri(uri string) (string, string) {
	if i := strings.Index(uri, "://"); i != -1 {
		return uri[0:i], uri[i+3:]
	}
	return "", uri
}

type errUnsupported struct {
	Source      string
	Destination string
}

func (e *errUnsupported) Error() string {
	return fmt.Sprintf("unsupported: %q => %q", e.Source, e.Destination)
}

func Sync(source string, destination string, parents bool, verbose bool) error {

	sourceScheme, sourcePath := splitUri(source)

	destinationScheme, destinationPath := splitUri(destination)

	if sourceScheme == "file" || sourceScheme == "" {

		if destinationScheme == "file" || destinationScheme == "" {
			return SyncLocalToLocal(sourcePath, destinationPath, parents, verbose)
		}

		if destinationScheme == "s3" {

			if verbose {
				fmt.Println("Uploading to AWS S3")
			}

			if err := s3util.CheckPath(destinationPath); err != nil {
				return fmt.Errorf("error with destination path %q: %w", destinationPath, err)
			}

			i := strings.Index(destinationPath, "/")

			if verbose {
				fmt.Println("creating AWS session")
			}

			s, err := awsutil.NewSession()
			if err != nil {
				return fmt.Errorf("error creating new session: %w", err)
			}

			return SyncLocalToS3(
				sourcePath,
				destinationPath[0:i],
				destinationPath[i+1:],
				s3manager.NewUploader(s),
				verbose,
			)
		}

		return &errUnsupported{Source: source, Destination: destination}
	}

	if sourceScheme == "s3" {

		if err := s3util.CheckPath(sourcePath); err != nil {
			return fmt.Errorf("error with source path %q: %w", sourcePath, err)
		}

		i := strings.Index(sourcePath, "/")

		if verbose {
			fmt.Println("creating AWS session")
		}

		if destinationScheme == "file" || destinationScheme == "" {

			if verbose {
				fmt.Println("Downloading from AWS S3")
			}

			s, err := awsutil.NewSession()
			if err != nil {
				return fmt.Errorf("error creating new session: %w", err)
			}

			return SyncS3ToLocal(
				sourcePath[0:i],
				sourcePath[i+1:],
				destinationPath,
				s3.New(s),
				s3manager.NewDownloader(s),
				verbose)
		}

		return &errUnsupported{Source: source, Destination: destination}
	}

	return &errUnsupported{Source: source, Destination: destination}

}

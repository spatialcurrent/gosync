// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/spatialcurrent/gosync/pkg/awsutil"
	"github.com/spatialcurrent/gosync/pkg/s3util"
)

const (
	WildcardCharacter = "*"

	MsgLocalFiles = "Synchronizing local directories"
	MsgNewSession = "Creating AWS Session"
	MsgUploadS3   = "Uploading to AWS S3"
	MsgDownloadS3 = "Downloading from AWS S3"
)

type errUnsupported struct {
	Source      string
	Destination string
}

func (e *errUnsupported) Error() string {
	return fmt.Sprintf("unsupported: %q => %q", e.Source, e.Destination)
}

type SyncInput struct {
	Source      string
	Destination string
	Parents     bool
	Limit       int
	Verbose     bool
	Credentials *credentials.Credentials
	PoolSize    int
	StopOnError bool
	Timeout     time.Duration
}

func Sync(input *SyncInput) error {

	sourceScheme, sourcePath := splitURI(input.Source)

	destinationScheme, destinationPath := splitURI(input.Destination)

	if sourceScheme == "file" || sourceScheme == "" {

		if destinationScheme == "file" || destinationScheme == "" {
			if input.Verbose {
				fmt.Println(MsgLocalFiles)
			}
			return SyncLocalToLocal(&SyncLocalToLocalInput{
				Source:      sourcePath,
				Destination: destinationPath,
				Parents:     input.Parents,
				Verbose:     input.Verbose,
			})
		}

		if destinationScheme == "s3" {

			if err := s3util.CheckPath(destinationPath); err != nil {
				return fmt.Errorf("error with destination path %q: %w", destinationPath, err)
			}

			i := strings.Index(destinationPath, "/")

			if input.Verbose {
				fmt.Println(MsgNewSession)
			}

			s, err := awsutil.NewSession(&awsutil.NewSessionInput{
				Credentials: input.Credentials,
				Verbose:     input.Verbose,
			})
			if err != nil {
				return fmt.Errorf("error creating new session: %w", err)
			}

			if input.Verbose {
				fmt.Println(MsgUploadS3)
			}

			return SyncLocalToS3(&SyncLocalToS3Input{
				Source:      sourcePath,
				Bucket:      destinationPath[0:i],
				KeyPrefix:   destinationPath[i+1:],
				Uploader:    s3manager.NewUploader(s),
				PoolSize:    input.PoolSize,
				Verbose:     input.Verbose,
				StopOnError: input.StopOnError,
				Timeout:     input.Timeout,
			})
		}

		return &errUnsupported{Source: input.Source, Destination: input.Destination}
	}

	if sourceScheme == "s3" {

		if err := s3util.CheckPath(sourcePath); err != nil {
			return fmt.Errorf("error with source path %q: %w", sourcePath, err)
		}

		i := strings.Index(sourcePath, "/")

		if destinationScheme == "file" || destinationScheme == "" {

			if input.Verbose {
				fmt.Println(MsgNewSession)
			}

			s, err := awsutil.NewSession(&awsutil.NewSessionInput{
				Credentials: input.Credentials,
				Verbose:     input.Verbose,
			})
			if err != nil {
				return fmt.Errorf("error creating new session: %w", err)
			}

			if input.Verbose {
				fmt.Println(MsgDownloadS3)
			}

			return SyncS3ToLocal(&SyncS3ToLocalInput{
				Bucket:      sourcePath[0:i],
				KeyPrefix:   sourcePath[i+1:],
				Destination: destinationPath,
				Client:      s3.New(s),
				Downloader:  s3manager.NewDownloader(s),
				Parents:     input.Parents,
				Limit:       input.Limit,
				Verbose:     input.Verbose,
				PoolSize:    input.PoolSize,
				StopOnError: input.StopOnError,
				Timeout:     input.Timeout,
			})
		}

		return &errUnsupported{Source: input.Source, Destination: input.Destination}
	}

	return &errUnsupported{Source: input.Source, Destination: input.Destination}

}

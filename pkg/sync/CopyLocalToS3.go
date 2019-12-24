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
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"strings"

	"os"
)

func CopyLocalToS3(source string, bucket string, key string, uploader *s3manager.Uploader) error {

	if strings.HasPrefix(source, "~") {
		return fmt.Errorf("source %q starts with \"~\"", source)
	}

	if strings.HasPrefix(bucket, "~") {
		return fmt.Errorf("destination bucket %q starts with \"~\"", bucket)
	}

	if strings.HasPrefix(key, "~") {
		return fmt.Errorf("destination key %q starts with \"~\"", key)
	}

	if uploader == nil {
		return fmt.Errorf("error uploading file to AWS S3: uploader is nil")
	}

	sourceFile, err := os.OpenFile(source, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("error opening source file at %q: %w", source, err)
	}
	defer sourceFile.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   sourceFile,
	})
	if err != nil {
		return fmt.Errorf("error uploading file to AWS S3: %w", err)
	}

	return nil
}

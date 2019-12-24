// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package s3util

import (
	"fmt"
	"strings"
)

type BucketError struct {
	Value string
}

func (e *BucketError) Error() string {
	return e.Value
}

func CheckBucket(bucket string) error {
	if strings.Contains(bucket, "~") {
		return &BucketError{Value: fmt.Sprintf("bucket cannot contain \"~\": %q", bucket)}
	}
	if strings.Contains(bucket, "*") {
		return &BucketError{Value: fmt.Sprintf("bucket cannot contain \"*\": %q", bucket)}
	}
	return nil
}

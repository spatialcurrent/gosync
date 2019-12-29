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

type PathError struct {
	Value string
}

func (e *PathError) Error() string {
	return e.Value
}

func CheckPath(path string) error {
	if !strings.Contains(path, "/") {
		return &PathError{Value: fmt.Sprintf("path is missing \"/\": %q", path)}
	}
	if strings.Contains(path, "~") {
		return &PathError{Value: fmt.Sprintf("path cannot contain \"~\": %q", path)}
	}
	if strings.Contains(path, "*") {
		return &PathError{Value: fmt.Sprintf("path cannot contain \"*\": %q", path)}
	}
	if strings.HasPrefix(path, "/") {
		return &PathError{Value: fmt.Sprintf("path cannot start with \"/\": %q", path)}
	}
	if strings.HasSuffix(path, "/") {
		return &PathError{Value: fmt.Sprintf("path cannot end with \"/\": %q", path)}
	}
	return nil
}

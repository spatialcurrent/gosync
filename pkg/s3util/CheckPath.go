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
	return nil
}

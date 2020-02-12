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

type KeyError struct {
	Value string
}

func (e *KeyError) Error() string {
	return e.Value
}

func CheckKey(prefix string) error {
	if strings.Contains(prefix, "~") {
		return &KeyError{Value: fmt.Sprintf("key cannot contain \"~\": %q", prefix)}
	}
	if strings.Contains(prefix, "*") {
		return &KeyError{Value: fmt.Sprintf("key cannot contain \"*\": %q", prefix)}
	}
	return nil
}

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

type KeyPrefixError struct {
	Value string
}

func (e *KeyPrefixError) Error() string {
	return e.Value
}

func CheckKeyPrefix(prefix string) error {
	if strings.Contains(prefix, "~") {
		return &KeyPrefixError{Value: fmt.Sprintf("key prefix cannot contain \"~\": %q", prefix)}
	}
	if strings.Contains(prefix, "*") {
		return &KeyPrefixError{Value: fmt.Sprintf("key prefix cannot contain \"*\": %q", prefix)}
	}
	return nil
}

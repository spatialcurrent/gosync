// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package s3util

type PathError struct {
	Value string
}

func (e *PathError) Error() string {
	return e.Value
}

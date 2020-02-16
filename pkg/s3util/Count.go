// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package s3util

import (
	"io"
)

// Count iterates through an iterator and returns the total number of objects.
func Count(it *Iterator) (int, error) {
	count := 0
	for {
		_, err := it.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return count, err
		}
		count += 1
	}
	return count, nil
}

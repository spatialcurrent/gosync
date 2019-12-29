// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

func maxLength(s []string) int {
	max := 0
	for _, v := range s {
		if i := len(v); i > max {
			max = i
		}
	}
	return max
}

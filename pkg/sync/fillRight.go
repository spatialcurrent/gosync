// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

import (
	"strings"
)

func fillRight(s string, l int) string {
	if len(s) < l {
		return s + strings.Repeat(" ", l-len(s))
	}
	return s
}

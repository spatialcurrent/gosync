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

func splitUri(uri string) (string, string) {
	if i := strings.Index(uri, "://"); i != -1 {
		return uri[0:i], uri[i+3:]
	}
	return "", uri
}

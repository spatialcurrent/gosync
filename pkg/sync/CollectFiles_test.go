// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectFiles(t *testing.T) {
	files, err := CollectFiles([]string{"testdata"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"testdata/a/b/doc.1.txt", "testdata/a/b/doc.2.txt", "testdata/a/doc.1.txt", "testdata/a/doc.2.txt", "testdata/b/doc.1.txt", "testdata/b/doc.2.txt", "testdata/doc.1.txt"}, files)
}

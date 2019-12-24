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

func TestSyncLocalToLocal(t *testing.T) {
	err := SyncLocalToLocal("testdata", "tmp", true, true)
	assert.NoError(t, err)
}

// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package cli

import (
	"github.com/spf13/pflag"
)

const (
	FlagParents = "parents"
	FlagLimit   = "limit"
)

// InitSyncFlags initializes the Sync flags.
func InitSyncFlags(flag *pflag.FlagSet) {
	flag.BoolP(FlagParents, "p", false, "create parent directories")
	flag.IntP(FlagLimit, "l", -1, "limit number of files copied")
}
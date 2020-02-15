// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package cli

import (
	"github.com/spf13/pflag"
)

const (
	FlagParents     = "parents"
	FlagLimit       = "limit"
	FlagPoolSize    = "pool-size"
	FlagStopOnError = "stop-on-error"

	DefaultPoolSize = 1
	DefaultLimit    = -1
)

// InitSyncFlags initializes the Sync flags.
func InitSyncFlags(flag *pflag.FlagSet) {
	flag.BoolP(FlagParents, "p", false, "create parent directories")
	flag.IntP(FlagLimit, "l", DefaultLimit, "limit number of files copied")
	flag.IntP(FlagPoolSize, "s", DefaultPoolSize, "pool size (number of concurrent downloads or uploads)")
	flag.BoolP(FlagStopOnError, "e", false, "stop on error")
}

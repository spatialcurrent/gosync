// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package cli

import (
	"github.com/spf13/pflag"
	"time"
)

const (
	FlagParents     = "parents"
	FlagLimit       = "limit"
	FlagPoolSize    = "pool-size"
	FlagStopOnError = "stop-on-error"
	FlagTimeout     = "timeout"

	DefaultPoolSize = 1
	DefaultLimit    = -1
)

// InitSyncFlags initializes the Sync flags.
func InitSyncFlags(flag *pflag.FlagSet) {
	flag.BoolP(FlagParents, "p", false, "create parent directories of destination if they do not exist")
	flag.IntP(FlagLimit, "l", DefaultLimit, "limit number of files copied")
	flag.IntP(FlagPoolSize, "s", DefaultPoolSize, "pool size (number of concurrent downloads or uploads)")
	flag.BoolP(FlagStopOnError, "e", false, "stop copying file if there is an error copying any of them")
	flag.DurationP(FlagTimeout, "t", 0*time.Second, "maximum duration for copying an individual file before aborting")
}

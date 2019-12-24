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
	FlagVerbose = "verbose"
)

// InitAWSFlags initializes the AWS flags.
func InitVerboseFlags(flag *pflag.FlagSet) {
	flag.BoolP(FlagVerbose, "v", false, "verbose output")
}

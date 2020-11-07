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
	FlagVersion = "version"
)

// InitVersionFlag initializes the version flag.
func InitVersionFlag(flag *pflag.FlagSet) {
	flag.Bool(FlagVersion, false, "show version")
}

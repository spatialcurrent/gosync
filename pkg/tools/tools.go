// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// +build tools

// This file exists to track tool dependencies. This is one of the recommended practices
// for handling tool dependencies in a Go module as outlined here:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	_ "github.com/gordonklaus/ineffassign"
	_ "github.com/kisielk/errcheck"
	_ "github.com/mitchellh/gox"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow"
	_ "honnef.co/go/tools/cmd/staticcheck"
)

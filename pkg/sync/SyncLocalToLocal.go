// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"path/filepath"
	"strings"

	"os"
)

type SyncLocalToLocalInput struct {
	Source      string
	Destination string
	Parents     bool
	Verbose     bool
}

func SyncLocalToLocal(input *SyncLocalToLocalInput) error {

	if strings.HasPrefix(input.Source, "~") {
		return fmt.Errorf("source %q starts with \"~\"", input.Source)
	}

	if strings.HasPrefix(input.Destination, "~") {
		return fmt.Errorf("destination %q starts with \"~\"", input.Destination)
	}

	sourcePaths, err := CollectFiles([]string{input.Source})
	if err != nil {
		return fmt.Errorf("error collecting files from %q: %w", input.Source, err)
	}

	if strings.Contains(input.Destination, WildcardCharacter) {
		return fmt.Errorf("destination cannot contain wildcard: %q", input.Destination)
	}

	destinationFile, err := os.Open(input.Destination)
	if err != nil {
		if os.IsNotExist(err) {
			if !input.Parents {
				return fmt.Errorf("destination directory %q does not exist and parents is not true", input.Destination)
			}
		} else {
			return fmt.Errorf("unable to open destination %q: %w", input.Destination, err)
		}
	} else {
		fileInfo, err := destinationFile.Stat()
		if err != nil {
			return fmt.Errorf("error stating destination %q: %w", input.Destination, err)
		} else {
			if !fileInfo.Mode().IsDir() {
				return fmt.Errorf("destination %q exists, but is not a directory", input.Destination)
			}
		}
	}

	sourceMaxLength := maxLength(sourcePaths)

	var g errgroup.Group
	for i, p := range sourcePaths {
		p := p
		r, err := filepath.Rel(input.Source, p)
		if err != nil {
			return fmt.Errorf("error calculating relative path between %q and %q: %w", input.Source, p, err)
		}
		destinationPath := filepath.Join(input.Destination, r)
		if input.Verbose {
			fmt.Printf("[ %d ] : %s    =>    %s\n", i+1, fillRight(p, sourceMaxLength), destinationPath)
		}
		g.Go(func() error {
			err := CopyLocalToLocal(p, destinationPath, input.Parents)
			if err != nil {
				return fmt.Errorf("error copying %q to %q: %w", p, destinationPath, err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	return nil
}

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

func SyncLocalToLocal(source string, destination string, parents bool, verbose bool) error {

	if strings.HasPrefix(source, "~") {
		return fmt.Errorf("source %q starts with \"~\"", source)
	}

	if strings.HasPrefix(destination, "~") {
		return fmt.Errorf("destination %q starts with \"~\"", destination)
	}

	sourcePaths, err := CollectFiles([]string{source})
	if err != nil {
		return fmt.Errorf("error collecting files from %q: %w", source, err)
	}

	if strings.Contains(destination, WildcardCharacter) {
		return fmt.Errorf("destination cannot contain wildcard: %q", destination)
	}

	destinationFile, err := os.Open(destination)
	if err != nil {
		if os.IsNotExist(err) {
			if !parents {
				return fmt.Errorf("destination directory %q does not exist and parents is not true", destination)
			}
		} else {
			return fmt.Errorf("unable to open destination %q: %w", destination, err)
		}
	} else {
		fileInfo, err := destinationFile.Stat()
		if err != nil {
			return fmt.Errorf("error stating destination %q: %w", destination, err)
		} else {
			if !fileInfo.Mode().IsDir() {
				return fmt.Errorf("destination %q exists, but is not a directory", destination)
			}
		}
	}

	sourceMaxLength := maxLength(sourcePaths)

	var g errgroup.Group
	for _, p := range sourcePaths {
		p := p
		r, err := filepath.Rel(source, p)
		if err != nil {
			return fmt.Errorf("error calculating relative path between %q and %q: %w", source, p, err)
		}
		destinationPath := filepath.Join(destination, r)
		if verbose {
			fmt.Println(fmt.Sprintf("%s    =>    %s", fillRight(p, sourceMaxLength), destinationPath))
		}
		g.Go(func() error {
			err := CopyLocalToLocal(p, destinationPath, parents)
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

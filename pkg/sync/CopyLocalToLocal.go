// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"os"
)

func CopyLocalToLocal(source string, destination string, parents bool) error {

	if strings.HasPrefix(source, "~") {
		return fmt.Errorf("source %q starts with \"~\"", source)
	}

	if strings.HasPrefix(destination, "~") {
		return fmt.Errorf("destination %q starts with \"~\"", destination)
	}

	parent := filepath.Dir(destination)
	if _, err := os.Stat(parent); err != nil {
		if os.IsNotExist(err) {
			if !parents {
				return fmt.Errorf("parent directory for destination %q does not exist and parents parameter is false", destination)
			}
			err := os.MkdirAll(parent, 0755)
			if err != nil {
				return fmt.Errorf("error creating parent directories for %q", destination)
			}
		} else {
			return fmt.Errorf("error stating destination parent %q", parent)
		}
	}

	sourceFile, err := os.OpenFile(source, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("error opening source file at %q: %w", source, err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("error creating destination file at %q: %w", source, err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("error copying from %q to %q", source, destination)
	}

	return nil
}

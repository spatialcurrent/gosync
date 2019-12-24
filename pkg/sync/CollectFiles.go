// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CollectFiles(in []string) ([]string, error) {
	out := make([]string, 0)
	for _, p := range in {
		if strings.Contains(p, WildcardCharacter) {
			matches, err := filepath.Glob(p)
			if err != nil {
				return out, fmt.Errorf("unable go glob source path %q: %w", p, err)
			}
			for _, match := range matches {
				err := filepath.Walk(match, func(p string, info os.FileInfo, err error) error {
					if err != nil {
						return fmt.Errorf("error walking %q: %w", p, err)
					}
					if info.Mode().IsRegular() {
						out = append(out, p)
					}
					return nil
				})
				if err != nil {
					return out, fmt.Errorf("error walking %q: %w", p, err)
				}
			}
		} else {
			err := filepath.Walk(p, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return fmt.Errorf("error walking %q: %w", p, err)
				}
				if info.Mode().IsRegular() {
					out = append(out, p)
				}
				return nil
			})
			if err != nil {
				return out, fmt.Errorf("error walking %q: %w", p, err)
			}
		}
	}
	return out, nil
}

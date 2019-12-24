// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package awsutil

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession() (*session.Session, error) {
	s, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("error creating new session: %w", err)
	}
	return s, nil
}

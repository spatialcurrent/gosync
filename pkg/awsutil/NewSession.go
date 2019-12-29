// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package awsutil

import (
	"errors"
	"fmt"
	
	"github.com/aws/aws-sdk-go/aws/session"
)

type NewSessionInput struct {
	Region  string
	Verbose bool
}

func NewSession(input *NewSessionInput) (*session.Session, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	config, err := NewConfig(&NewConfigInput{
		Region:  input.Region,
		Verbose: input.Verbose,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating new config: %w", err)
	}
	s, err := session.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("error creating new session: %w", err)
	}
	return s, nil
}

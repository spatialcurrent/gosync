// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package awsutil

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
)

type NewConfigInput struct {
	Region  string
	Verbose bool
}

func NewConfig(input *NewConfigInput) (*aws.Config, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	config := &aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(input.Verbose),
	}

	if region := input.Region; len(region) > 0 {
		config.Region = aws.String(region)
	}

	return config, nil
}

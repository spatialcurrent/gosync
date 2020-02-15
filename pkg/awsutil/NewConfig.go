// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package awsutil

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type NewConfigInput struct {
	Credentials *credentials.Credentials
	Region      string
	Verbose     bool
}

func NewConfig(input *NewConfigInput) (*aws.Config, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	config := &aws.Config{
		Credentials:                   input.Credentials,
		CredentialsChainVerboseErrors: aws.Bool(input.Verbose),
	}

	if region := input.Region; len(region) > 0 {
		config.Region = aws.String(region)
	}

	return config, nil
}

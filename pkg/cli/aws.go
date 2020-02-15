// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package cli

import (
	"github.com/spf13/pflag"
)

const (
	FlagAWSDefaultRegion = "aws-default-region"
	FlagAWSMFASerial     = "aws-mfa-serial"
	FlagAWSRegion        = "aws-region"
	FlagAWSRoleARN       = "aws-role-arn"
)

// InitAWSFlags initializes the AWS flags.
func InitAWSFlags(flag *pflag.FlagSet) {
	flag.String(FlagAWSDefaultRegion, "", "AWS Default Region")
	flag.String(FlagAWSMFASerial, "", "AWS MFA Serial")
	flag.StringP(FlagAWSRegion, "", "", "AWS Region (overrides default region)")
	flag.String(FlagAWSRoleARN, "", "AWS Role ARN")
}

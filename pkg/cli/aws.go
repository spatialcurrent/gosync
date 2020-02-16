// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package cli

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"time"
)

const (
	FlagAWSDefaultRegion      = "aws-default-region"
	FlagAWSMFASerial          = "aws-mfa-serial"
	FlagAWSRegion             = "aws-region"
	FlagAWSRoleARN            = "aws-role-arn"
	FlagAWSAssumeRoleDuration = "aws-assume-role-duration"

	DefaultAWSAssumeRoleDuration = 15 * time.Minute

	MinimumAWSAssumeRoleDuration = 15 * time.Minute
)

// InitAWSFlags initializes the AWS flags.
func InitAWSFlags(flag *pflag.FlagSet) {
	flag.String(FlagAWSDefaultRegion, "", "AWS Default Region")
	flag.String(FlagAWSMFASerial, "", "AWS MFA Serial")
	flag.StringP(FlagAWSRegion, "", "", "AWS Region (overrides default region)")
	flag.String(FlagAWSRoleARN, "", "AWS Role ARN")
	flag.Duration(FlagAWSAssumeRoleDuration, DefaultAWSAssumeRoleDuration, "Expiry duration of the STS credentials for assuming a role.")
}

func CheckAWSConfig(v *viper.Viper) error {
	assumeRoleDuration := v.GetDuration(FlagAWSAssumeRoleDuration)
	if assumeRoleDuration < 15*time.Minute {
		return fmt.Errorf("%q value %q is invalid, expecting value greater than or equal to %q", FlagAWSAssumeRoleDuration, assumeRoleDuration, MinimumAWSAssumeRoleDuration)
	}
	return nil
}

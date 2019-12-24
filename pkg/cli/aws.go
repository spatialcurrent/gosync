// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	FlagAWSProfile                         = "aws-profile"
	FlagAWSDefaultRegion                   = "aws-default-region"
	FlagAWSRegion                          = "aws-region"
	FlagAWSAccessKeyId                     = "aws-access-key-id"
	FlagAWSSecretAccessKey                 = "aws-secret-access-key"
	FlagAWSSessionToken                    = "aws-session-token"
	FlagAWSSecurityToken                   = "aws-security-token"
	FlagAWSContainerCredentialsRelativeUri = "aws-container-credentials-relative-uri"
)

// InitAWSFlags initializes the AWS flags.
func InitAWSFlags(flag *pflag.FlagSet) {
	flag.StringP(FlagAWSProfile, "", "", "AWS Profile")
	flag.StringP(FlagAWSDefaultRegion, "", "", "AWS Default Region")
	flag.StringP(FlagAWSRegion, "", "", "AWS Region")
	flag.StringP(FlagAWSAccessKeyId, "", "", "AWS Access Key ID")
	flag.StringP(FlagAWSSecretAccessKey, "", "", "AWS Secret Access Key")
	flag.StringP(FlagAWSSessionToken, "", "", "AWS Session Token")
	flag.StringP(FlagAWSSecurityToken, "", "", "AWS Security Token")
	flag.StringP(FlagAWSContainerCredentialsRelativeUri, "", "", "AWS Container Credentials Relative URI")
}

// CheckAWSConfig checks the AWS configuration.
func CheckAWSConfig(v *viper.Viper) error {
	return nil
}

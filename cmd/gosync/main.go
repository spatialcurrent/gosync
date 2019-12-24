// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/spatialcurrent/gosync/pkg/cli"
	"github.com/spatialcurrent/gosync/pkg/sync"
)

func initFlags(flag *pflag.FlagSet) {
	cli.InitAWSFlags(flag)
	cli.InitSyncFlags(flag)
	cli.InitVerboseFlags(flag)
}

func initViper(cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()
	err := v.BindPFlags(cmd.Flags())
	if err != nil {
		return v, fmt.Errorf("error binding flag set to viper: %w", err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv() // set environment variables to overwrite config
	return v, nil
}

func checkConfig(v *viper.Viper, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expecting 2 positional argument (source and destination), but found %d arguments", len(args))
	}
	return nil
}

func main() {
	cmd := &cobra.Command{
		Use:                   "gosync [SOURCE] [DESTINATION]",
		DisableFlagsInUseLine: true,
		Short:                 "gosync",
		Long:                  `gosyc is a super simple command line program for synchronizing two directories, including support for AWS S3 buckets.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			v, err := initViper(cmd)
			if err != nil {
				return fmt.Errorf("error initializing viper: %w", err)
			}

			if errConfig := checkConfig(v, args); errConfig != nil {
				return errConfig
			}

			verbose := v.GetBool(cli.FlagVerbose)
			parents := v.GetBool(cli.FlagParents)

			source := args[0]
			destination := args[1]

			err = sync.Sync(source, destination, parents, verbose)
			if err != nil {
				return fmt.Errorf("error syncing from %q to %q: %w", source, destination, err)
			}

			return nil

		},
	}
	initFlags(cmd.Flags())

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "gosync: "+err.Error())
		fmt.Fprintln(os.Stderr, "Try gosync --help for more information.")
		os.Exit(1)
	}
}

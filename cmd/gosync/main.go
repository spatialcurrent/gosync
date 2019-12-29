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
		return fmt.Errorf("expecting 2 positional arguments for source and destination, but found %d arguments", len(args))
	}
	return nil
}

func main() {
	cmd := &cobra.Command{
		Use:                   "gosync SOURCE DESTINATION",
		DisableFlagsInUseLine: true,
		Short:                 "gosync",
		Long:                  `gosyc is a super simple command line program for synchronizing two directories specified by URI.  gosync currently supports local directories and AWS S3 buckets as a source or destination.  AWS S3 buckets are specified using the "s3://" scheme.  Local files are specified using the "file://" scheme or a path without a scheme.  gosync synchronizes regular files and will create directories as needed if the parents flag is set.`,
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {

			v, err := initViper(cmd)
			if err != nil {
				return fmt.Errorf("error initializing viper: %w", err)
			}

			if errConfig := checkConfig(v, args); errConfig != nil {
				return errConfig
			}

			verbose := v.GetBool(cli.FlagVerbose)

			source := args[0]
			destination := args[1]

			err = sync.Sync(&sync.SyncInput{
				Source:      source,
				Destination: destination,
				Parents:     v.GetBool(cli.FlagParents),
				Limit:       v.GetInt(cli.FlagLimit),
				Verbose:     verbose,
			})
			if err != nil {
				return fmt.Errorf("error syncing from %q to %q: %w", source, destination, err)
			}

			if verbose {
				fmt.Println("Done.")
			}

			return nil

		},
	}
	initFlags(cmd.Flags())

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "gosync: "+err.Error())
		fmt.Fprintln(os.Stderr, "Try \"gosync --help\" for more information.")
		os.Exit(1)
	}
}

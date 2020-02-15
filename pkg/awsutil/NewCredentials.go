// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package awsutil

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/spatialcurrent/goprompt/pkg/prompt"
)

const (
	MaxNumberOfMFATokenRequests = 100
)

type NewCredentialsInput struct {
	Session      *session.Session
	Role         string
	SerialNumber string
}

func NewCredentials(input *NewCredentialsInput) *credentials.Credentials {
	if len(input.SerialNumber) > 0 {
		count := 0
		mutex := &sync.Mutex{}
		return stscreds.NewCredentials(input.Session, input.Role, func(p *stscreds.AssumeRoleProvider) {
			p.SerialNumber = aws.String(input.SerialNumber)
			p.TokenProvider = func() (string, error) {
				mutex.Lock()
				defer mutex.Unlock()
				if count > MaxNumberOfMFATokenRequests {
					return "", fmt.Errorf("too many MFA token requests, exceeds limit of %d", MaxNumberOfMFATokenRequests)
				}
				v, err := prompt.String("gosync: enter MFA token", false, true) // will loop on blank entries
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				count += 1
				return v, err
			}
		})
	}
	return stscreds.NewCredentials(input.Session, input.Role)
}

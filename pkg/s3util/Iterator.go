// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package s3util

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
)

type Iterator struct {
	client  *s3.S3
	bucket  string
	prefix  string
	marker  *string
	objects []*s3.Object
	err     error
	eof     bool
}

func (it *Iterator) Error() error {
	return it.err
}

func (it *Iterator) Next() (*s3.Object, error) {
	if it.eof {
		return nil, io.EOF
	}
	if len(it.objects) == 0 {
		listObjectsOutput, err := it.client.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(it.bucket),
			Prefix: aws.String(it.prefix),
			Marker: it.marker,
		})
		if err != nil {
			it.err = fmt.Errorf("error listing objets from marker %q", aws.StringValue(it.marker))
			return nil, it.err
		} else {
			it.err = nil
		}
		if len(listObjectsOutput.Contents) == 0 {
			it.eof = true
			return nil, io.EOF
		}
		it.objects = listObjectsOutput.Contents
		it.marker = listObjectsOutput.Marker
	}

	object := it.objects[0]
	it.objects = it.objects[1:]
	return object, nil
}

type NewIteratorInput struct {
	Client *s3.S3
	Bucket string
	Prefix string
}

func NewIterator(input *NewIteratorInput) *Iterator {
	return &Iterator{
		client:  input.Client,
		bucket:  input.Bucket,
		prefix:  input.Prefix,
		marker:  nil,
		objects: make([]*s3.Object, 0),
		eof:     false,
	}
}

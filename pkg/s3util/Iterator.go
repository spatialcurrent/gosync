// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package s3util

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Iterator struct {
	client  *s3.S3
	bucket  string
	prefix  string
	maxKeys int
	marker  *string
	objects []*s3.Object
	err     error
	last    bool
	eof     bool
}

func (it *Iterator) Reset(input *NewIteratorInput) {
	it.client = input.Client
	it.bucket = input.Bucket
	it.prefix = input.Prefix
	it.maxKeys = input.MaxKeys
	it.marker = nil
	it.objects = it.objects[:0]
	it.last = false
	it.eof = false
}

func (it *Iterator) Error() error {
	return it.err
}

func (it *Iterator) Next() (*s3.Object, error) {
	if it.eof {
		return nil, io.EOF
	}
	if len(it.objects) == 0 {
		// if the last page was already requested
		if it.last {
			it.eof = true
			return nil, io.EOF
		}
		listObjectsOutput, err := it.client.ListObjects(&s3.ListObjectsInput{
			Bucket:  aws.String(it.bucket),
			Prefix:  aws.String(it.prefix),
			Marker:  it.marker,
			MaxKeys: aws.Int64(int64(it.maxKeys)),
		})
		if err != nil {
			it.err = fmt.Errorf("error listing objets from marker %q: %w", aws.StringValue(it.marker), err)
			return nil, it.err
		} else {
			it.err = nil
		}
		if len(listObjectsOutput.Contents) == 0 {
			it.eof = true
			return nil, io.EOF
		}
		it.objects = listObjectsOutput.Contents
		// if results are not truncated, then you've returned all the results
		if !aws.BoolValue(listObjectsOutput.IsTruncated) {
			it.last = true
		} else {
			// set the marker to the key of the last object returned
			it.marker = listObjectsOutput.Contents[len(listObjectsOutput.Contents)-1].Key
		}
	}

	object := it.objects[0]
	it.objects = it.objects[1:]
	return object, nil
}

type NewIteratorInput struct {
	Client  *s3.S3
	Bucket  string
	Prefix  string
	MaxKeys int
}

func NewIterator(input *NewIteratorInput) *Iterator {
	return &Iterator{
		client:  input.Client,
		bucket:  input.Bucket,
		prefix:  input.Prefix,
		maxKeys: input.MaxKeys,
		marker:  nil,
		objects: make([]*s3.Object, 0),
		last:    false,
		eof:     false,
	}
}

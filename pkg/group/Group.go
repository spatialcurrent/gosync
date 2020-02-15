// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package group contains a struct that wraps errgroup.Group.
//	- "golang.org/x/sync/errgroup"
package group

import (
	"sync"

	"golang.org/x/sync/errgroup"
)

type Group struct {
	errgroup.Group
	pool        chan bool
	stopOnError bool
	limit       int
	count       int
	stop        bool
	mutex       *sync.Mutex
}

func (g *Group) Go(f func() error) {
	// if stopped, then return immediately
	g.mutex.Lock()
	if g.stop {
		g.mutex.Unlock()
		return
	}
	g.mutex.Unlock()
	// if limit was exceeded, then return immediately without scheduling goroutine
	if g.limit > 0 && g.count >= g.limit {
		return
	}
	g.Group.Go(func() error {
		// if stopped, then return immediately
		g.mutex.Lock()
		if g.stop {
			g.mutex.Unlock()
			return nil
		}
		g.mutex.Unlock()
		// allocate file descriptor
		g.pool <- true
		// unallocate file descriptor after function executes
		defer func() { <-g.pool }()
		// execute given function in this goroutine
		err := f()
		// if function returned an error and stopOnError is set, then set stop.
		if err != nil && g.stopOnError {
			g.stop = true
		}
		return err
	})
	g.count += 1
}

func New(poolSize int, limit int, stopOnError bool) *Group {
	return &Group{
		pool:        make(chan bool, poolSize),
		stopOnError: stopOnError,
		stop:        false,
		limit:       limit,
		count:       0,
		mutex:       &sync.Mutex{},
	}
}

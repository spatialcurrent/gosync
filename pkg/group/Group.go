// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

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
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.stop {
		return
	}
	if g.limit > 0 && g.count >= g.limit {
		return
	}
	g.Group.Go(func() error {
		g.mutex.Lock()
		defer g.mutex.Unlock()
		if g.stop {
			return nil
		}
		g.pool <- true
		defer func() { <-g.pool }()
		err := f()
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

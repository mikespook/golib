// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package iptpool

import (
	"sync"
)

var (
	DefaultMaxIdle = 32
)

type CreateFunc func() ScriptIpt
type EventFunc func(ScriptIpt) error

type ScriptIpt interface {
    Exec(name string, params interface{}) error
    Init(path string) error
    Final() error
	Bind(name string, item interface{}) error
}

type IptPool struct {
	maxIdle int
	mu sync.Mutex
	freeIpt []ScriptIpt
	create CreateFunc

	OnCreate EventFunc
}

func NewIptPool(create CreateFunc) (pool *IptPool) {
	pool = &IptPool{
		maxIdle: DefaultMaxIdle,
		create: create,
	}
	return
}

func (pool *IptPool) Get() (ipt ScriptIpt) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if n := len(pool.freeIpt); n > 0 {
		ipt = pool.freeIpt[n-1]
		pool.freeIpt = pool.freeIpt[:n-1]
		return
	}
	ipt = pool.New()
	return
}

func (pool *IptPool) Put(ipt ScriptIpt) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if n := len(pool.freeIpt); n > pool.maxIdle {
		return
	}
	pool.freeIpt = append(pool.freeIpt, ipt)
	return
}

func (pool *IptPool) GetMaxIdle() int {
	return pool.maxIdle
}

func (pool *IptPool) SetMaxIdle(maxIdle int) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if maxIdle <= 0 {
		maxIdle = 0
	}
	pool.maxIdle = maxIdle
	if maxIdle == 0 {
		pool.freeIpt = nil
		return
	}
	if n := len(pool.freeIpt); n > maxIdle {
		pool.freeIpt = pool.freeIpt[:maxIdle-1]
	}
}

func (pool *IptPool) New() ScriptIpt {
	ipt := pool.create()
	if pool.OnCreate != nil {
		pool.OnCreate(ipt)
	}
	return ipt
}

func (pool *IptPool) Free() map[int]error {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	emap := make(map[int]error)
	for k, ipt := range pool.freeIpt {
		if err := ipt.Final(); err != nil {
			emap[k] = err
		}
	}
	return emap
}

func (pool *IptPool) Length() int {
	return len(pool.freeIpt)
}

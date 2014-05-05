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
	p sync.Pool
}

func NewIptPool(create CreateFunc) (pool *IptPool) {
	f := func() interface{} {
		return create()
	}
	return &IptPool{sync.Pool{New: f}}
}

func (pool *IptPool) Get() (ipt ScriptIpt) {
	return pool.p.Get().(ScriptIpt)
}

func (pool *IptPool) Put(ipt ScriptIpt) {
	pool.p.Put(ipt)
}

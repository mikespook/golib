// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package scriptpool

import (
	"sync"
)

var (
	DefaultMaxIdle = 32
)

type NewFunc func() ScriptIpt

type ScriptIpt interface {
    Exec(name string, params interface{}) error
    Init(path string, pool ScriptPool) error
    Final() error
	Bind(name string, item interface{}) error
}

type ScriptPool struct {
	maxIdle int
	mu sync.Mutex
	freeIpt []ScriptIpt
	newFunc NewFunc
}

func NewScriptPool(newFunc NewFunc, preassign bool) (pool *ScriptPool) {
	pool = &ScriptPool{
		maxIdle: DefaultMaxIdle,
		newFunc: newFunc,
	}
	if preassign {
		for i := 0; i < pool.maxIdle; i ++ {
			pool.Put(pool.newFunc())
		}
	}
	return
}

func (pool *ScriptPool) Get() (ipt ScriptIpt) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if n := len(pool.freeIpt); n > 0 {
		ipt = pool.freeIpt[n-1]
		pool.freeIpt = pool.freeIpt[:n-1]
		return
	}
	ipt = pool.newFunc()
	return
}

func (pool *ScriptPool) Put(ipt ScriptIpt) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if n := len(pool.freeIpt); n > pool.maxIdle {
		return
	}
	pool.freeIpt = append(pool.freeIpt, ipt)
	return
}


func (pool *ScriptPool) SetMaxIdle(maxIdle int) {
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

func (pool *ScriptPool) Free() map[int]error {
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

func (pool *ScriptPool) Length() int {
	return len(pool.freeIpt)
}

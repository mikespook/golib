// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package iptpool

import (
    "testing"
)

type testIpt struct {}

func (t *testIpt) Exec(name string, params interface{}) error {return nil}
func (t *testIpt) Init(path string, pool *IptPool) error {return nil}
func (t *testIpt) Final() error {return nil}
func (t *testIpt) Bind(name string, item interface{}) error {return nil}

func newTestIpt() ScriptIpt {
	return &testIpt{}
}

func TestPool(t *testing.T) {
	pool := NewIptPool(newTestIpt)
	if n := pool.Length(); n != 0 {
		t.Error("Wrong pool length: %d", n)
	}
	for i := 0; i < DefaultMaxIdle; i ++ {
		ipt1 := pool.Get()
		ipt2 := pool.Get()
		pool.Put(ipt1)
		pool.Put(ipt2)
	}
	if n := pool.Length(); n != 2 {
		t.Error("Wrong pool length: %d", n)
	}
}

func TestPoolPreassign(t *testing.T) {
	pool := NewIptPool(newTestIpt)
	for i := 0; i < pool.GetMaxIdle(); i ++ {
		pool.Put(pool.New())
	}
	if n := pool.Length(); n != DefaultMaxIdle {
		t.Error("Wrong pool length: %d", n)
	}
	for i := 0; i < DefaultMaxIdle; i ++ {
		ipt1 := pool.Get()
		ipt2 := pool.Get()
		pool.Put(ipt1)
		pool.Put(ipt2)
	}
	if n := pool.Length(); n != DefaultMaxIdle {
		t.Error("Wrong pool length: %d", n)
	}
}

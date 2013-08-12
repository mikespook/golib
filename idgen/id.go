// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package idgen

import (
    "labix.org/v2/mgo/bson"
    "github.com/mikespook/golib/autoinc"
)

type IdGen interface {
    Id() interface{}
}

// ObjectId
type objectId struct {}

func (id *objectId) Id() interface{} {
    return bson.NewObjectId().Hex()
}

func NewObjectId() IdGen {
    return &objectId{}
}

// AutoIncId
type autoincId struct {
    *autoinc.AutoInc
}

func (id *autoincId) Id() interface{} {
    return id.AutoInc.Id()
}

func NewAutoIncId() IdGen {
    return &autoincId{autoinc.New(1, 1)}
}

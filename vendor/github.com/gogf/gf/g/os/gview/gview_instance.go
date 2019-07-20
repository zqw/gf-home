// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

import "github.com/gogf/gf/g/container/gmap"

const (
	// Default group name for instance usage.
	DEFAULT_INSTANCE_NAME = "default"
)

var (
	// Instances map.
	instances = gmap.NewStrAnyMap()
)

// Instance returns an instance of View with default settings.
// The parameter <name> is the name for the instance.
func Instance(name ...string) *View {
	key := DEFAULT_INSTANCE_NAME
	if len(name) > 0 {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		return New()
	}).(*View)
}

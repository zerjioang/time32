//
// Created by zerjioang
// https://github/zerjioang
// Copyright (c) 2020. All rights reserved.
//
// SPDX-License-Identifier: GPL-3.0
//

package time32

import (
	"sync/atomic"
	"time"
)

// lastTime stores the epoch value of last time reading
var lastTime atomic.Value
var lastUnix atomic.Value
var lastUnixNano atomic.Value

func init() {
	// store initial value
	tt := time.Now()
	lastTime.Store(tt)
	lastUnix.Store(tt.Unix())
	lastUnixNano.Store(tt.UnixNano())

	// run each 0.1 seconds (aka precision)
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				lastTime.Store(t)
				lastUnix.Store(t.Unix())
				lastUnixNano.Store(t.UnixNano())
			}
		}
	}()
}

// ReuseTime is a function that reuses last readed epoch value
// this function is meant to be used on high demanding applications that require
// time value readings with high frequency. Instead of making a syscall on every request,
// last time value is cached. Cache duration has a window of 0.1s so all calls requested during
// that period will reuse the same epoch time value
func ReuseTime() time.Time {
	return lastTime.Load().(time.Time)
}

func ReuseUnix() int64 {
	return lastUnix.Load().(int64)
}

func ReuseUnixNano() int64 {
	return lastUnixNano.Load().(int64)
}
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

// lastEpoch stores the epoch value of last time reading
var lastEpoch atomic.Value

func init() {
	// store initial value
	lastEpoch.Store(int64(Epoch()))

	// run each 0.1 seconds (aka precision)
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				lastEpoch.Store(t.Unix())
			}
		}
	}()
}

// ReuseEpoch is a function that reuses last readed epoch value
// this function is meant to be used on high demanding applications that require
// time value readings with high frequency. Instead of making a syscall on every request,
// last time value is cached. Cache duration has a window of 0.1s so all calls requested during
// that period will reuse the same epoch time value
func ReuseEpoch() int64 {
	return lastEpoch.Load().(int64)
}

//
// Created by zerjioang
// https://github/zerjioang
// Copyright (c) 2020. All rights reserved.
//
// SPDX-License-Identifier: GPL-3.0
//

// Package time32 implements a time.Now() features of Go language in order to fetch Epoch time our UTC time
// without the need of using internal pointers for UTC location data ( *loc )
package time32

/*
Time32 Defines our own time unit which will always hold epoch time
in millis. Example: 1588228661
4294967295
int64 size: 8 bytes
uint32 size: 4 bytes

x2 size reduction changing data size
*/
type Time32 uint32

func (t Time32) AddDate(days int) Time32 {
	v := int(t) + (days * 86400)
	return Time32(v)
}

func (t *Time32) setTime(now uint32) {
	*t = Time32(now)
}

// Epoch Returns current server epoch millis time without
// GC dealing with *loc pointers
func Epoch() Time32 {
	return Time32(get_now())
}

// Epoch Returns current server epoch millis time without
// GC dealing with *loc pointers
func get_now() uint32 {
	sec, nsec, mono := time_now()
	mono -= startNano
	var wall uint64
	var ext int64
	var unsec = uint64(nsec)
	sec += unixToInternal - minWall
	var usec = uint64(sec)
	if usec>>33 != 0 {
		wall, ext = unsec, sec+minWall
	} else {
		wall, ext = hasMonotonic|usec<<nsecShift|unsec, mono
	}
	tsec := ext
	if wall&hasMonotonic != 0 {
		tsec = wallToInternal + int64(wall<<1>>(nsecShift+1))
	}
	return uint32(tsec + internalToUnix)
}

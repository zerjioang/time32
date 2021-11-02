//
// Created by zerjioang
// https://github/zerjioang
// Copyright (c) 2020. All rights reserved.
//
// SPDX-License-Identifier: GPL-3.0
//

package time32

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func ExampleEpoch() {
	fmt.Println(Epoch())
}

func TestNow(t *testing.T) {
	t.Run("custom-time", func(t *testing.T) {
		tt := Now()
		r := reflect.ValueOf(tt)
		t.Log(r)
		t.Log(binary.Size(tt))
	})
	t.Run("standard-go-time", func(t *testing.T) {
		tt := time.Now()
		r := reflect.ValueOf(tt)
		t.Log(r)
		t.Log(binary.Size(tt))
	})
	t.Run("custom-epoch", func(t *testing.T) {
		tt := Epoch()
		t.Log(binary.Size(tt))
		t.Log(tt)
	})
	t.Run("standard-go-epoch", func(t *testing.T) {
		tt := time.Now().Unix()
		t.Log(binary.Size(tt))
		t.Log(tt)
	})
	t.Run("both-epoch", func(t *testing.T) {
		tt := time.Now().Unix()
		t.Log(binary.Size(tt))
		t.Log(tt)
		t2 := Epoch()
		t.Log(binary.Size(t2))
		t.Log(t2)
	})
}

func BenchmarkNow(b *testing.B) {
	// BenchmarkNow/epoch-custom-12         	     232	   5111623 ns/op	   0.00 MB/s	       0 B/op	       0 allocs/op
	b.Run("epoch-custom", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(1)
		b.ResetTimer()
		var stamps [100000]Time32
		for i := 0; i < b.N; i++ {
			for i := 0; i < 100000; i++ {
				stamps[i] = Epoch()
			}
		}
	})
	//BenchmarkNow/epoch-standard-go-12         	     249	   4805626 ns/op	   0.00 MB/s	       0 B/op	       0 allocs/op
	b.Run("epoch-standard-go", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(1)
		b.ResetTimer()
		var stamps [100000]int64
		for i := 0; i < b.N; i++ {
			for i := 0; i < 100000; i++ {
				stamps[i] = time.Now().Unix()
			}
		}
	})
	b.Run("custom", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Now()
		}
	})
	b.Run("standard-go-time", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Now()
		}
	})
	b.Run("custom-ref", func(b *testing.B) {
		// make a benchmark in where compiler
		// optimizations do not remove our variable
		var tt Time
		b.ReportAllocs()
		b.SetBytes(1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tt = Now()
		}
		if tt.IsZero() {
			b.Log("time is zero")
		}
	})
	b.Run("standard-go-time-ref", func(b *testing.B) {
		// make a benchmark in where compiler
		// optimizations do not remove our variable
		var tt time.Time
		b.ReportAllocs()
		b.SetBytes(1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tt = time.Now()
		}
		if tt.IsZero() {
			b.Log("time is zero")
		}
	})
	b.Run("reuse-epoch", func(b *testing.B) {
		// make a benchmark in where compiler
		// optimizations do not remove our variable
		b.ReportAllocs()
		b.SetBytes(1)
		b.ResetTimer()
		var ep int64
		for i := 0; i < b.N; i++ {
			ep = ReuseEpoch()
		}
		if ep == 0 {
			b.Log("time is zero")
		}
	})
}

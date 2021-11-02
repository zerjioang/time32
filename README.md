# `uint32` time for Golang

> Note: this package is based on `Golang` official time implementation.

## Features

* Get current time as `uint32` epoch (millis).
* Same performance as `time` package from standard Go library.
* Less GC pressure by removing timezone location data related fields.
* Time cache for high speed time requests. Default update frequency: 0.1s

## Differences from `time.Time` package

* `*loc` pointer usage has been removed from `Time` struct to avoid pointer pressure on GC cycles.
* Included a method `Epoch()` that returns current epoch time as `uint32` instead of `int64`. This means, we can store our time data in **4 bytes**.

## Usage

```go

package main

import (
  "fmt"
  "github.com/zerjioang/time32"
)

func main(){
  tt := time32.Epoch()
  fmt.Println(tt)
}
```

```bash
# prints current time in epoch format
Example: 1601452800
```

An example is provided at ./example/main.go

## Cached timing

For application that require a constant access to system time, a cached mechanism is included that keeps last time on a **0.1s** refresh rate.
To use this caching features, available methods are:

* `ReuseTime`
* `ReuseUnix`
* `ReuseUnixNano`

Previous method will return last value within a 0.1s window. Note that this feature might be useful for adding a timestamp to logs, expiration check, etc.

## Performance

```bash
goos: linux
goarch: amd64
pkg: github.com/zerjioang/time32

BenchmarkNow/epoch-custom-12         	     356	       3314537 ns/op	     0.00 MB/s	       1 B/op	       0 allocs/op
BenchmarkNow/epoch-standard-go-12    	     345	       3301546 ns/op	     0.00 MB/s	       1 B/op	       0 allocs/op
BenchmarkNow/custom-12               	37442151	         33.07 ns/op	    30.24 MB/s	       0 B/op	       0 allocs/op
BenchmarkNow/standard-go-time-12     	38039060	         32.16 ns/op	    31.09 MB/s	       0 B/op	       0 allocs/op
BenchmarkNow/custom-ref-12           	36875410	         32.14 ns/op	    31.11 MB/s	       0 B/op	       0 allocs/op
BenchmarkNow/standard-go-time-ref-12 	36313864	         32.98 ns/op	    30.32 MB/s	       0 B/op	       0 allocs/op
BenchmarkNow/reuse-time-12             932550987	         1.270 ns/op	   787.44 MB/s	       0 B/op	       0 allocs/op
BenchmarkNow/reuse-unix-12             920748625	         1.271 ns/op	   786.86 MB/s	       0 B/op	       0 allocs/op
BenchmarkNow/reuse-unixnano-12         888746485	         1.309 ns/op	   764.14 MB/s	       0 B/op	       0 allocs/op
BenchmarkNow/reuse-time-unixnano-12    529924074	         2.292 ns/op	   436.22 MB/s	       0 B/op	       0 allocs/op
```

## License

All rights reserved to **@zerjioang** under **GNU GPL v3** license

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

 * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 * Uses GPLv3 license described below

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

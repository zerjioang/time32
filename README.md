# `uint32` time for Golang

> Note: this package is based on `Golang` official time implementation.

## Features

* Get current time as `uint32` epoch (millis).
* Same performance as `time` package from standard Go library.
* Less GC pressure by removing timezone location data related fields.

## Differences from `time.Time` package

* `*loc` pointer usage has been removed from `Time` struct to avoid pointer pressure on GC cycles.
* Included a method `Epoch()` that returns current epoch time as `uint32` instead of `int64`. This means, we can store our time data in **4 bytes**.

## TL;DR - usage

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
# prints current time in seconds
Example: 1601452800
```

## License

All rights reserved to **@zerjioang** under **GNU GPL v3** license

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

 * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 * Uses GPLv3 license described below

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

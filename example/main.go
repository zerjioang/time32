//
// Created by zerjioang
// https://github/zerjioang
// Copyright (c) 2020. All rights reserved.
//
// SPDX-License-Identifier: GPL-3.0
//

package main

import (
	"fmt"
	"github.com/zerjioang/time32"
	"time"
)

func main() {
	fmt.Println("current time using time32.Epoch is: ", time32.Epoch())
	fmt.Println("current time using time32.ReuseEpoch is: ", time32.ReuseEpoch())
	fmt.Println("current time using time.Now().Unix() is: ", time.Now().Unix())
}

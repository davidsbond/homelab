// Copyright 2012 Aaron Jacobs. All Rights Reserved.
// Author: aaronjjacobs@gmail.com (Aaron Jacobs)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build 386 arm,!arm64 mips mipsle

package cmac

import (
	"log"
	"unsafe"
)

// XOR the blockSize bytes starting at a and b, writing the result over dst.
func xorBlock(
	dstPtr unsafe.Pointer,
	aPtr unsafe.Pointer,
	bPtr unsafe.Pointer) {
	// Check assumptions. (These are compile-time constants, so this should
	// compile out.)
	const wordSize = unsafe.Sizeof(uintptr(0))
	if blockSize != 4*wordSize {
		log.Panicf("%d %d", blockSize, wordSize)
	}

	// Convert.
	a := (*[4]uintptr)(aPtr)
	b := (*[4]uintptr)(bPtr)
	dst := (*[4]uintptr)(dstPtr)

	// Compute.
	dst[0] = a[0] ^ b[0]
	dst[1] = a[1] ^ b[1]
	dst[2] = a[2] ^ b[2]
	dst[3] = a[3] ^ b[3]
}

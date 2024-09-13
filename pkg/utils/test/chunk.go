package test

// min returns a if it's smaller than b, otherwise b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Chunk chunks a given slice into sub-slices of up to n elements of s.
// All but the last sub-slice will have size n.
// All sub-slices are clipped to have no capacity beyond the length.
// If s is empty, the sequence is empty: there is no empty slice in the sequence.
// Chunk panics if n is less than 1.
//
// Note: This is slightly adjusted and backported to Go 1.20.
//
// FIXME: Starting with Go 1.23, the actual [slices.Chunk] function should be used instead.
//
// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Source: https://cs.opensource.google/go/go/+/refs/tags/go1.23.0:src/slices/iter.go
func Chunk[Slice ~[]E, E any](s Slice, n int) []Slice {
	var res []Slice
	for i := 0; i < len(s); i += n {
		// Clamp the last chunk to the slice bound as necessary.
		end := min(n, len(s[i:]))

		// Set the capacity of each chunk so that appending to a chunk does
		// not modify the original slice.
		res = append(res, s[i:i+end:i+end])
	}
	return res
}

package utils

import (
	_ "container/heap"
)

// IntHeap ...
type IntHeap []int32

// Len sort interface in heap
func (h IntHeap) Len() int {
	return len(h)
}

// Less sort interface in heap
func (h IntHeap) Less(ia, ib int) bool {
	if ia < 0 || ib < 0 || ia >= len(h) || ib >= len(h) {
		return false
	}
	return h[ia] < h[ib]
}

// Less sort interface in heap
func (h IntHeap) Swap(ia, ib int) {
	if ia < 0 || ib < 0 || ia >= len(h) || ib >= len(h) {
		return
	}
	h[ia], h[ib] = h[ib], h[ia]
}

// Push interface of heap
func (h *IntHeap) Push(x interface{}) {
	if x == nil {
		return
	}
	// Push and Pop use pointer receivers because they modify the slice's length, not just its contents.
	*h = append(*h, x.(int32))
}

// Pop interface of heap
func (h *IntHeap) Pop() interface{} {
	old := *h
	len := len(*h)
	if len < 1 {
		return nil
	}
	x := old[len-1]
	*h = old[0 : len-1]
	return x
}

// Data ...
func (h *IntHeap) Data() []int32 {
	rst := make([]int32, len(*h))
	copy(rst, *h)
	return rst
}

package main

import "errors"

type MinHeap struct {
	size     int
	top      int
	data     []*HuffmanTreeNode
	capacity int
}

func NewMinHeap(capacity int) *MinHeap {
	if capacity < 0 {
		// default capacity: 2
		capacity = 2
	}
	return &MinHeap{
		size:     0,
		top:      -1,
		data:     make([]*HuffmanTreeNode, capacity),
		capacity: capacity,
	}
}

func (h *MinHeap) Insert(x *HuffmanTreeNode) {
	if h.IsFull() {
		buf := make([]*HuffmanTreeNode, h.capacity*2)
		copy(buf, h.data)
		h.data = buf
		h.capacity *= 2
	}
	h.top++
	h.data[h.top] = x
	h.size++
	h.upNode(h.top)
	return
}

func (h *MinHeap) Delete() (x *HuffmanTreeNode, err error) {
	if h.IsEmpty() {
		return nil, errors.New("Heap is empty")
	}
	if h.size < h.capacity/2 {
		buf := make([]*HuffmanTreeNode, h.size+1)
		copy(buf, h.data)
		h.data = buf
		h.capacity = h.size + 1
	}
	x = h.data[0]
	h.data[0] = h.data[h.top]
	h.data[h.top] = nil
	h.top--
	h.size--
	h.downNode(0)
	return x, nil
}

func (h *MinHeap) BuildHeap() {
	for i := h.size / 2; i >= 0; i-- {
		h.downNode(i)
	}
	return
}

func (h *MinHeap) IsEmpty() bool {
	return h.size == 0
}

func (h *MinHeap) IsFull() bool {
	return h.size == h.capacity
}

func (h *MinHeap) upNode(i int) {
	child, parent := i, i/2
	for child >= 0 {
		if child+1 < h.size && h.data[child+1].curVal < h.data[child].curVal {
			child++
		}
		if h.data[child].curVal >= h.data[parent].curVal {
			break
		}
		h.data[child], h.data[parent] = h.data[parent], h.data[child]
		child, parent = parent, parent/2
	}
}

func (h *MinHeap) downNode(i int) {
	child, parent := i*2, i
	for child < h.size {
		if child+1 < h.size && h.data[child+1].curVal < h.data[child].curVal {
			child++
		}
		if h.data[child].curVal >= h.data[parent].curVal {
			break
		}
		h.data[child], h.data[parent] = h.data[parent], h.data[child]
		child, parent = child*2, child
	}
}

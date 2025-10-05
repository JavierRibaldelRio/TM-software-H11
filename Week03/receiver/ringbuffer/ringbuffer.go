package ringbuffer

import (
	"fmt"
	"sync"
)

type RingBuffer[T any] struct {
	mu   sync.RWMutex
	head *ringNode[T]
	size int
}

type ringNode[T any] struct {
	ele  T
	next *ringNode[T]
}

// Creates the ring
func NewRing[T any](size int) *RingBuffer[T] {

	// Checks if the Ring can be created
	if size < 2 {
		panic("The size of the Ring Buffer must be atleast 2")
	}

	// First
	firstNode := &ringNode[T]{}

	// Before node
	pre := firstNode

	// Inserts n-1 nodes
	for i := 0; i < size-1; i++ {

		n := &ringNode[T]{}

		pre.next = n

		pre = n

	}

	// Last node points to the first

	pre.next = firstNode

	return &RingBuffer[T]{head: firstNode, size: size}

}

// Removes the oldest and adds a new one
func (ring *RingBuffer[T]) Add(v T) {

	fmt.Print(ring.Read())

	// Blocks to write
	ring.mu.Lock()

	// Overwrites
	ring.head.ele = v

	// Rotation
	ring.head = ring.head.next

	ring.mu.Unlock()
}

// Returns the size of the ring

func (ring *RingBuffer[T]) Size() int {

	return ring.size
}

// Returns the content of the ring from old to new
func (ring *RingBuffer[T]) Read() []T {

	buffer := make([]T, 0, ring.size)

	ring.mu.RLock()
	defer ring.mu.RUnlock()

	for n := ring.head; ; n = n.next {
		buffer = append(buffer, n.ele)
		if n.next == ring.head {
			break
		}
	}
	return buffer

}

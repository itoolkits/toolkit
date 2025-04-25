// hash set

package collect

type HashSet[T comparable] struct {
	container map[T]struct{}
	nilVal    T // do not init it
}

// NewHashSet create hash set
func NewHashSet[T comparable]() *HashSet[T] {
	return &HashSet[T]{
		container: map[T]struct{}{},
	}
}

// NewHashSetBySlice - create hash set by slice
func NewHashSetBySlice[T comparable](arr []T) *HashSet[T] {
	h := NewHashSet[T]()
	h.Add(arr...)
	return h
}

// NewHashSetByHashSet - create hash set by hash set
func NewHashSetByHashSet[T comparable](hs ...*HashSet[T]) *HashSet[T] {
	h := NewHashSet[T]()
	if len(hs) < 1 {
		return h
	}
	for _, ele := range hs {
		ele.Range(func(t T) bool {
			h.Add(t)
			return true
		})
	}
	return h
}

// Add - add element
func (h *HashSet[T]) Add(args ...T) {
	for i := range args {
		// in most cases, init value is meaningless
		if args[i] == h.nilVal {
			continue
		}
		h.container[args[i]] = nilStructObj
	}
}

// Size - container size
func (h *HashSet[T]) Size() int {
	return len(h.container)
}

// Contains - contains ele
func (h *HashSet[T]) Contains(ele T) bool {
	_, ok := h.container[ele]
	return ok
}

// Remove - remove
func (h *HashSet[T]) Remove(ele T) {
	delete(h.container, ele)
}

// Clear - clear container
func (h *HashSet[T]) Clear() {
	clear(h.container)
}

// Range - loop element
func (h *HashSet[T]) Range(fn func(T) bool) {
	for ele := range h.container {
		fnR := fn(ele)
		if !fnR {
			break
		}
	}
}

// ToSlice - to slice
func (h *HashSet[T]) ToSlice() []T {
	rst := make([]T, 0, len(h.container))
	for k := range h.container {
		rst = append(rst, k)
	}
	return rst
}

// AddHashSet - add hash set
func (h *HashSet[T]) AddHashSet(hs ...*HashSet[T]) {
	if len(hs) < 1 {
		return
	}
	for i := range hs {
		hs[i].Range(func(ele T) bool {
			h.Add(ele)
			return true
		})
	}
}

// RemoveHashSet - remove hash set
func (h *HashSet[T]) RemoveHashSet(hs ...*HashSet[T]) {
	if len(hs) < 1 {
		return
	}
	for i := range hs {
		hs[i].Range(func(ele T) bool {
			h.Remove(ele)
			return true
		})
	}
}

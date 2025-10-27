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

// NewHashSetByEle create hash set by element
func NewHashSetByEle[T comparable](args ...T) *HashSet[T] {
	c := &HashSet[T]{
		container: map[T]struct{}{},
	}
	c.Add(args...)
	return c
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
func (h *HashSet[T]) Add(args ...T) *HashSet[T] {
	for i := range args {
		// in most cases, init value is meaningless
		if args[i] == h.nilVal {
			continue
		}
		h.container[args[i]] = nilStructObj
	}
	return h
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
func (h *HashSet[T]) Remove(ele T) *HashSet[T] {
	delete(h.container, ele)
	return h
}

// Clear - clear container
func (h *HashSet[T]) Clear() *HashSet[T] {
	clear(h.container)
	return h
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
func (h *HashSet[T]) AddHashSet(hs ...*HashSet[T]) *HashSet[T] {
	if len(hs) < 1 {
		return h
	}
	for i := range hs {
		hs[i].Range(func(ele T) bool {
			h.Add(ele)
			return true
		})
	}
	return h
}

// RemoveHashSet - remove hash set
func (h *HashSet[T]) RemoveHashSet(hs ...*HashSet[T]) *HashSet[T] {
	if len(hs) < 1 {
		return h
	}
	for i := range hs {
		hs[i].Range(func(ele T) bool {
			h.Remove(ele)
			return true
		})
	}
	return h
}

// Intersection get intersection set
func (h *HashSet[T]) Intersection(target *HashSet[T]) *HashSet[T] {
	if target == nil || target.Size() < 1 || h.Size() < 1 {
		return nil
	}

	a := h
	b := target
	if a.Size() > b.Size() {
		a = target
		b = h
	}

	rst := NewHashSet[T]()
	a.Range(func(ele T) bool {
		if b.Contains(ele) {
			rst.Add(ele)
		}
		return true
	})

	return rst
}

// Union get Union set
func (h *HashSet[T]) Union(target *HashSet[T]) *HashSet[T] {
	rst := NewHashSet[T]()
	rst.AddHashSet(target, h)
	return rst
}

// Except get Except by set
func (h *HashSet[T]) Except(target *HashSet[T]) *HashSet[T] {
	rst := NewHashSet[T]()
	h.Range(func(ele T) bool {
		if target != nil && !target.Contains(ele) {
			rst.Add(ele)
		}
		return true
	})
	return rst
}

// CompareHashSet compare hashset
func CompareHashSet[T comparable](a, b *HashSet[T]) (*HashSet[T], *HashSet[T]) {
	return NewHashSetByHashSet(a).RemoveHashSet(b),
		NewHashSetByHashSet(b).RemoveHashSet(a)
}

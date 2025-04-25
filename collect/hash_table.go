// table struct, use hash map implement
// R - row
// C - column
// V - value

//    c1,c2,c3
// r1 1, 2, 3
// r2 4, 5, 6

package collect

type HashTable[R, C comparable, V any] struct {
	container map[R]map[C]V
	nilVal    V // do not init it
}

// NewHashTable - create hash table
func NewHashTable[R, C comparable, V any]() *HashTable[R, C, V] {
	return &HashTable[R, C, V]{
		container: make(map[R]map[C]V),
	}
}

// Add - add element
func (t *HashTable[R, C, V]) Add(r R, c C, v V) {
	columns, h := t.container[r]
	if !h {
		t.container[r] = map[C]V{
			c: v,
		}
	} else {
		columns[c] = v
	}
}

// Get - get element, by row and column
func (t *HashTable[R, C, V]) Get(r R, c C) (V, bool) {
	columns, h := t.container[r]
	if !h {
		v, h := columns[c]
		return v, h
	}
	return t.nilVal, false
}

// Delete - delete element by row and column
func (t *HashTable[R, C, V]) Delete(r R, c C) {
	columns, h := t.container[r]
	if !h {
		return
	}
	_, h = columns[c]
	if !h {
		return
	}
	delete(columns, c)
}

// Clear - clear container
func (t *HashTable[R, C, V]) Clear() {
	clear(t.container)
}

// Range - loop element
func (t *HashTable[R, C, V]) Range(fn func(R, C, V) bool) {
	for r, columns := range t.container {
		for c, ele := range columns {
			fnR := fn(r, c, ele)
			if !fnR {
				break
			}
		}
	}
}

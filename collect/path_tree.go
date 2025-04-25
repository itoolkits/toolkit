package collect

type PathTreeNode[K comparable, V any] struct {
	v        V
	children map[K]PathTreeNode[K, V]
}
type PathTree[K comparable, V any] struct {
	root *PathTreeNode[K, V]
}

func NewPathTree[K comparable, V any]() *PathTree[K, V] {
	return &PathTree[K, V]{
		root: &PathTreeNode[K, V]{
			children: make(map[K]PathTreeNode[K, V]),
		},
	}
}

func (t *PathTree[K, V]) Add(ks []K, v V) {

}

func (t *PathTree[K, V]) Get(ks []K) {

}

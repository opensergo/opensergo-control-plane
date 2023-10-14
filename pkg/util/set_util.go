package util

type Set[T comparable] map[T]struct{}

type String = Set[string]

func NewWithLength[T comparable](l int) Set[T] {
	return make(Set[T], l)
}

// New creates a new Set with the given items.
func New[T comparable](items ...T) Set[T] {
	s := NewWithLength[T](len(items))
	return s.InsertAll(items...)
}

// InsertAll adds the items to this Set.
func (s Set[T]) InsertAll(items ...T) Set[T] {
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

func (s Set[T]) Difference(s2 Set[T]) Set[T] {
	result := New[T]()
	for key := range s {
		if !s2.Contains(key) {
			result.Insert(key)
		}
	}
	return result
}

// Contains returns whether the given item is in the set.
func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

// Insert a single item to this Set.
func (s Set[T]) Insert(item T) Set[T] {
	s[item] = struct{}{}
	return s
}

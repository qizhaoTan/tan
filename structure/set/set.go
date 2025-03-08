package set

// Set implementation using generic.
type Set[T comparable] map[T]struct{}

// New gives a new Set object.
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// ToSlice set to a slice
func (s Set[T]) ToSlice() []T {
	arr := make([]T, 0, len(s))
	for i := range s {
		arr = append(arr, i)
	}
	return arr
}

// Insert a val to the set
func (s Set[T]) Insert(val T) {
	s[val] = struct{}{}
}

// Find the val is in the set
func (s Set[T]) Find(val T) bool {
	_, ok := s[val]
	return ok
}

// Remove the val from the set
func (s Set[T]) Remove(val T) {
	delete(s, val)
}

// Size of the set
func (s Set[T]) Size() int {
	return len(s)
}

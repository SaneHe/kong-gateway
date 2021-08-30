package main

var itemExists = struct{}{}

type Set struct {
	items map[string]struct{}
}

func NewSet(values ...string) *Set {
	set := &Set{items: make(map[string]struct{})}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

func (set *Set) Add(items ...string) {
	for _, item := range items {
		set.items[item] = itemExists
	}
}

func (set *Set) Contains(items ...string) bool {
	for _, item := range items {
		if _, contains := set.items[item]; !contains {
			return false
		}
	}
	return true
}

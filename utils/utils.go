package utils

import (
	"sort"
)

// Remove all elements in exp from src
func Except(src []string, exp []string) []string {
	sort.Strings(exp)
	remove_indexes := make([]int, 0)
	for i, v := range src {
		s_i := sort.Search(len(exp), func(i int) bool { return exp[i] >= v })
		if s_i < len(exp) && exp[s_i] == v {
			remove_indexes = append(remove_indexes, i)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(remove_indexes)))
	for _, ri := range remove_indexes {
		src = remove(src, ri)
	}
	return src
}

// Remove element with index i from slice
// (unsave: don't check boundaries of index)
func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
